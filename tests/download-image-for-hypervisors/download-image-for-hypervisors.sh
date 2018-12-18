# download image script
# Script used for checking that downloading an image for each hypervisor works

cutOffFileSize="10000"
moth="{{.MothershipBinPathAndName}}"

# Setup: Build an image to download
instAlias={{.OriginalAlias}}
instID=$($moth inspect-instance $instAlias -o id)
naclID=$($moth push-nacl tests/download-image-for-hypervisors/interface.nacl {{.BuilderID}} -o id)
imgID=$($moth build Starbase {{.BuilderID}} --instance $instID --nacl $naclID --waitAndPrint)

# Download the image for hypervisors qemu, vcloud and virtualbox
for hypervisor in vcloud virtualbox qemu
do
    sent=$[$sent + 1]
    downloadedImgName="hypervisor-test-image"
    # Download image
    raw+=$($moth pull-image $imgID $downloadedImgName --format $hypervisor 2>&1)
    if [ "$?" -ne "0" ]; then
      raw+="download of image for hypervisor $hypervisor failed"
      continue
    fi

    # Verify image size
    img_file_size=$(du -sk $downloadedImgName | awk '{print $1}')
    if [ "$img_file_size" -lt "$cutOffFileSize" ]; then
      raw+="$downloadedImgName size: $img_file_size, is too small"
      continue
    fi

    # Verify that image contains GRUB
    strings $downloadedImgName | grep GRUB 2>&1 > /dev/null
    if [ "$?" -ne "0" ]; then
      raw+="image for hypervisor $hypervisor does not contain GRUB"
      continue
    fi

    # Verify correct file format
    case "$hypervisor" in
      "vcloud" ) file $downloadedImgName | grep "VMware4 disk image" 2>&1 > /dev/null;;
      "virtualbox" ) file $downloadedImgName | grep "VirtualBox Disk Image" 2>&1 > /dev/null;;
      "qemu" ) file $downloadedImgName | grep "DOS/MBR boot sector" 2>&1 > /dev/null;;
    esac
    if [ "$?" -ne "0" ]; then
      raw+="image for hypervisor $hypervisor does not seem to be correct file format"
      raw+=$(file $downloadedImgName 2>&1)
      continue
    fi

    received=$[$received + 1]
done

# If none of the commands above failed it means that we were successful
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
