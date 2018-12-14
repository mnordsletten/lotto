# create issue script
# Script used for checking that it is possible to create an issue
set -e

moth="{{.MothershipBinPathAndName}}"

sent=0
received=0

# create-issue --name <name> --type <type> --description <description>
# delete-issue
# inspect-issue
# issues
# issuetypes
# pull-issue

issueNameBase="lotto-issue"

# Create an issue 3 times and verify that the issue was created
for i in {1..3}
do
    issueName="$issueNameBase-$i"
    sent=$[$sent + 1]
    # Create issue
    createdIssueID=$($moth create-issue --name $issueName --type Deployment --description "This is an issue created by Lotto" --waitAndPrint)
    # Verify that the issue was actually created
    nameOfIssueCreated=$($moth inspect-issue $createdIssueID -o json | jq -r '.name')
    if [[ "$nameOfIssueCreated" == "$issueName" ]]; then
        received=$[$received + 1]
    fi
done

if [ "$sent" -eq "$received" ]; then
  success=true
fi

if [ -z $success ]; then success=false; fi
if [ -z $sent ]; then sent=0; fi
if [ -z $received ]; then received=0; fi
if [ -z $rate ]; then rate=0; fi
if [ -z $raw ]; then raw=""; fi
jq \
  --argjson success $success \
  --argjson sent $sent \
  --argjson received $received \
  --argjson rate $rate \
  --arg raw "$raw" \
  '. |
  .["success"]=$success |
  .["sent"]=$sent |
  .["received"]=$received |
  .["rate"]=$rate |
  .["raw"]=$raw
  '<<<'{}'
