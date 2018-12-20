# alias update script
# Script used for checking that changing the alias of the instance works

moth="{{.MothershipBinPathAndName}}"
instAlias={{.OriginalAlias}}
ID=$($moth inspect-instance $instAlias -o id)
newAlias=$(echo alias-test-"$(date | shasum | cut -d " " -f 1)")

alias=""
sent=0
received=0

# Update alias 50 times (alternate between instAlias and newAlias)
for i in {1..50}
do
    # Set alias to use
    if (( $i % 2 )); then
        alias=$instAlias
    else
        alias=$newAlias
    fi

    sent=$[$sent + 1]
    # Change alias
    raw+=$($moth instance-alias $ID $alias 2>&1)
    # Verify that the alias was actually changed
    existingAlias=$($moth inspect-instance $ID -o json | jq -r '.alias')
    if [[ "$existingAlias" == "$alias" ]]; then
        received=$[$received + 1]
    fi
    sleep 0.2
done

# Reset to original alias
raw+=$($moth instance-alias $ID $instAlias 2>&1)

# All commands above succeeded
if [ "$sent" -eq "$received" ]; then
  success=true
fi

jq \
  --argjson success ${success:-false} \
  --argjson sent ${sent:-0} \
  --argjson received ${received:-0} \
  --argjson rate ${rate:-0} \
  --arg raw "${raw:-''}" \
  '. |
  .["success"]=$success |
  .["sent"]=$sent |
  .["received"]=$received |
  .["rate"]=$rate |
  .["raw"]=$raw
  '<<<'{}'
