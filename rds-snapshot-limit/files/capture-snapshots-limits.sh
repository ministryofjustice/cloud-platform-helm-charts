#!/bin/sh

kubectl get ns | awk 'FNR == 1 {next} {print $1}' > all-namespaces.txt # Dump all namespaces to file to iterate through later (skip the fisrst line as its header)

# Iterate through each namespace and capture the snapshot limit

echo "DB_NAME | SNAPSHOT_LIMIT" >> dblist-with-limit.txt

while read NAMESPACE; do

  SNAPSHOT_LIMIT=$(kubectl get ns $NAMESPACE -o jsonpath="{.metadata.annotations.max-manual-db-snapshots}") #Get the max limit number for each namespace if it exists
  if [ ! -z "$SNAPSHOT_LIMIT" ] # If not null then this namespace has a max snapshot annotation
  then
        #From the namespace we can now get the rds secret name
        #RDS_SECRET_NAME=$(kubectl get serviceaccount ns-access -o jsonpath='{.secrets[0].name}')
        RDS_SECRET_NAME=$(kubectl get secret -n $NAMESPACE | awk '/rds/{print $1}')
        #Now that we have the secret name we can get the db_name (which is a key in the secret)
        DB_NAME=$(kubectl get secret $RDS_SECRET_NAME -o jsonpath="{.data.rds_server}" -n $NAMESPACE | base64 --decode)
        echo $DB_NAME"|"$SNAPSHOT_LIMIT >> dblist-with-limit.txt #Finally output both the db_name and snapshot_limit to a text file to be used 
  fi
done <all-namespaces.txt







