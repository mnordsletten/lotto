# build and deploy script
# Script used for checking that building an image and deploying it to an instance (live update) works
set -e

moth="{{.MothershipBinPathAndName}}"
instAlias={{.OriginalAlias}}
instID=$($moth inspect-instance $instAlias -o id)
naclID=$($moth push-nacl tests/build-and-deploy/interface.nacl {{.BuilderID}} -o id)
tagBase="image"

# Build and deploy 10 times
for i in {1..10}
do
    tag="$tagBase-$i"
    sent=$[$sent + 1]
    # Build
    imgID=$($moth build Starbase {{.BuilderID}} --instance $instID --nacl $naclID --tag $tag --waitAndPrint)
    # Deploy
    raw+=$($moth deploy $instID $imgID --wait 2>&1)
    # Check if the instance now runs the built image and that it reports the given tag:
    cmdOutImgID=$($moth inspect-instance $instID -o json | jq -r '.imageId')
    cmdOutTag=$($moth inspect-instance $instID -o json | jq -r '.tag')
    if [[ "$cmdOutImgID" == *"$imgID"* ]] && [[ "$cmdOutTag" == *"$tag"* ]]; then
        received=$[$received + 1]
    fi
done

# If none of the commands above failed it means that we were successful
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
