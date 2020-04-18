#!/bin/bash
FILES=./*.jpg
for f in $FILES
do
   g=$(echo $f | sed 's/.jpg/-clean.jpg/')
   echo "Clean $f to $g"
   textcleaner $f $g
done
