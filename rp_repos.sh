#!/bin/bash

PRJS_NUM=$1
TAGS_NUM=$2
PULL_NUM=$3

PROJECT_NAME=temp
HARBOR_ADDR=11.11.11.12

# color setting
L1="\033[42;30;1m"
L2="\033[44;30;1m"
L3="\033[46;30;1m"
END="\033[0m"


help() {
    echo "Usage: ./import.sh <project_count> <tags_count> <pull_count>"
    exit 0
}

retention_policy_test() {

    echo
    echo "====== [[ prepare ]] ======"
    echo

	# 0. prepare docker images
    echo -e "$L1 --> prepare docker images $END"
    echo
	docker pull hello-world
	#docker pull busybox
    echo

	sleep 2

    # 1. create $PRJS_NUM projects for test
    echo -e "$L1 --> create $PRJS_NUM projects for test $END"
    echo
    for (( i=1; i<=$PRJS_NUM; i++ ))
    do
        # create both public and private projects
        if (( i % 2 == 0 )) ; then
            PUBLIC=1
        else
            PUBLIC=0
        fi

        echo
        #echo -e "$L2 ----> create: project_name=${PROJECT_NAME}_$i   public=${PUBLIC} $END"
        echo -e "$L2 ----> create: project_name=${PROJECT_NAME}_$i   public=1 $END"
        echo
        #./harbor-go-client prj_create --project_name="${PROJECT_NAME}_$i" --public=${PUBLIC}
        ./harbor-go-client prj_create --project_name="${PROJECT_NAME}_$i" --public=1

		rands=$((RANDOM % $TAGS_NUM))
        echo
        echo -e "$L2 ----> generate a random number ($rands) of tags and push onto Harbor $END"
        echo

		for r in $(seq $rands)
		do
            # 2. create tags, push repos and tags
			docker tag hello-world:latest ${HARBOR_ADDR}/${PROJECT_NAME}_${i}/hello-world:v${r}
			docker push ${HARBOR_ADDR}/${PROJECT_NAME}_${i}/hello-world:v${r} >/dev/null && echo "docker push ${HARBOR_ADDR}/${PROJECT_NAME}_${i}/hello-world:v${r}" || echo "Failed"

    		# 3. make a random number of pull
            loop=$((RANDOM % $PULL_NUM))
        	echo -e "$L3 ----> pull ${PROJECT_NAME}_${i}/hello-world:v${r} random times ($loop) $END"
			echo
            for j in $(seq $loop)
            do
			    docker pull ${HARBOR_ADDR}/${PROJECT_NAME}_${i}/hello-world:v${r} >/dev/null
            done
		done
    done

    # 4. do Retention Policy analysis
    echo
    echo "====== [[ retention policy analysis ]] ======"
    echo

    # 5. delete specified repos under projects
	./harbor-go-client rp_repos

	sleep 2

    # 6. delete all projects created for test (only projects whitout repos and policies can be deleted)
    echo
    echo "====== [[ clean ]] ======"
    echo

    echo
    echo -e "$L1 --> delete all projects created for test $END"
    echo

    IDS=$(./harbor-go-client prjs_list --name=${PROJECT_NAME} | grep "project_id" | awk '{print $2}' | sed -r 's/,//g')
	for id in $IDS
	do
    	./harbor-go-client prj_del --project_id=${id}
        echo
	done

    # 7. delete local docker tags and images
    echo -e "$L1 --> delete all tags and images created locally $END"
    echo
    docker rmi -f $(docker images hello-world -q)
}

doLogin() {
    echo
    echo -e "$L1 --> login $END"
    echo
    ./harbor-go-client login -u admin -p Harbor12345
    echo

    echo "----- docker login -----"
    docker login --username admin --password Harbor12345 $HARBOR_ADDR
    echo

}

doLogout() {
    echo
    echo -e "$L1 --> logout $END"
    echo
    ./harbor-go-client logout
    echo

    echo "----- docker logout -----"
    docker logout $HARBOR_ADDR
    echo
}


# ---------------

if [ $# -ne 3 ] ; then
    help
fi

doLogin
retention_policy_test
doLogout




