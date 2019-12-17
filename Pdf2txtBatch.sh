#!/bin/bash

# Assumes pdftotext is installed: http://www.xpdfreader.com/about.html

FILES='repace/with/path/to/pdfs'            #~/Documents/Research/Papers/*.pdf
for f in $FILES
do
 echo "Processing file .."
 pdftotext -enc UTF-8 $f
done
