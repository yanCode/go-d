#!/bin/bash
folders_to_remove=(
  "3000"
  "4000"
)

 remove_folder(){
   local folder="$1_network"
   if [ -d "$folder" ]; then
     echo "Removing folder: $folder"
     rm -rf "$folder"
     if [ $? -ne 0 ]; then
       echo "Failed to remove folder: $folder"
     fi
   else
     echo "Folder $folder does not exist, skipping it..."
   fi
 }
 echo "starting to clean network folders...."

 for folder in "${folders_to_remove[@]}"; do
   remove_folder "$folder"
 done