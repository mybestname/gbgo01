#!/usr/bin/env bash

# Download files from the url inputs, overwrite the same name local files.
# files :
#    https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
#    https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto
#    https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/httpbody.proto
# usages :
#    $ api=https://raw.githubusercontent.com/googleapis/googleapis/master/google/api
#    $ ./down.bash ${api}/annotations.proto ${api}http.proto ${api}httpbody.proto

for url in "$@"
do
  #echo "$url"
  file=$(basename $url)
  DOWN="wget --quiet --no-check-certificate --content-disposition -O $file"
  # TODO: if wget not find use curl
  # DOWN="curl -LJO"
  printf "\033[30;1mDownload\033[0m $file"
  $DOWN $url
  if [ $? -eq 0 ]; then
  printf "\033[32;1m OK\033[0m\n"
  else
  printf "\033[31;1m Failed\033[0m\n"
  fi
done

