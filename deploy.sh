#!/bin/bash
export ANSIBLE_HOST_KEY_CHECKING=False
ansible-playbook -vv -i hosts ansible.yaml -e "wxAppId=$wxAppId wxappSecret=$wxappSecret listenAddr=$listenAddr mysqlDSN=$mysqlDSN image_tag=$IMAGE_TAG"

