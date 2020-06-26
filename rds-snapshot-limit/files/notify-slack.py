#!/usr/bin/python3
import boto3
import subprocess
import os
import requests
import json

client = boto3.client('rds')

def get_business_unit(rds_server):
    arn = "arn:aws:rds:eu-west-2:754256621582:db:"+rds_server
    rdstags = client.list_tags_for_resource(ResourceName=arn)
    for tag in rdstags['TagList']:
      if tag['Key'] == 'business-unit':
        return tag['Value']

def slack_notification(business_unit,rds_server,actual_snapshot_count,expected_snapshot_limit):

    webhook_url = os.environ['SLACK_HOOK_URL']
   # slack_data = {'text': 'Team: '+business_unit+' | RDS Server: '+rds_server+' has exceeded its manual snapshot limit of '+str(expected_snapshot_limit)+ '. Actual number of manual snapshots is '+str(actual_snapshot_count)}
    
    
    slack_data = {
        "text": "RDS Instance: Manual snapshots exceeded limit",
         "attachments": [
            {
            "fallback": "Plain-text summary of the attachment.",
            "color": "#ECB22E",
            "author_name": "RDS Intance Name: "+rds_server,
            "author_link": "https://eu-west-2.console.aws.amazon.com/rds/home?region=eu-west-2#database:id="+rds_server+";is-cluster=false;tab=maintenance-and-backups",
            "text": "Number of manual snapshots ("+str(actual_snapshot_count)+") exceeds the limit of ("+str(expected_snapshot_limit)+")",
            "fields": [
                {
                    "title": "Priority",
                    "value": "Medium"
                }
            ],
            "footer": "Team: "+business_unit
        }
    ]
}        
    
    response = requests.post(webhook_url, data=json.dumps(slack_data),headers={'Content-Type': 'application/json'})
    if response.status_code != 200:
        raise ValueError(
            'Request to slack returned an error %s, the response is:\n%s'
            % (response.status_code, response.text)
        )


subprocess.call(['/opt/capture-snapshots-limits.sh'])

dblist_with_limit_file = open("dblist-with-limit.txt")

next(dblist_with_limit_file)
for line in dblist_with_limit_file:
    dblist = line.rstrip().split("|")
    rds_server = dblist[0]
    expected_snapshot_limit = int(dblist[1])
    response = client.describe_db_snapshots(DBInstanceIdentifier=rds_server,SnapshotType='manual')
    actual_snapshot_count = len(response["DBSnapshots"]) # Get the number of actual manual snapshots
    print('DB: '+rds_server+ '- expected limit: '+str(expected_snapshot_limit)+' - actual snapshot count: '+ str(actual_snapshot_count))
    if actual_snapshot_count > expected_snapshot_limit: # If actual number of snapshots exceeds limit then send slack notification

        slack_notification(get_business_unit(rds_server),rds_server,actual_snapshot_count,expected_snapshot_limit)
dblist_with_limit_file.close()


