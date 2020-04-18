# gradex-testset

## Why?

I'm writing code to automate aspects of processing exam scripts, with around 500 - 3000 pages per exam. Many are likely to come in at around 1500 pages.

This repo is an attempt to create a fairly large set of pages, with unique handwriting on each page, and combine them into a set of pdf files with a normal distribution of page lengths around a mean of 15 pages.

It's for rehearsal and familiarity purposes. Many features are developed with small sample files, so this is for building confidence that we've at least attempted to flush out any scaling issues.


## Image sources

The free-to-register 
[IAM handwriting database](http://www.fki.inf.unibe.ch/databases/iam-handwriting-database) just happens to contain slightly over 1500 pages of documents with handwriting. Perfect.

The site uses basic auth, so you can download the archives using wget (strangle chrome/linux was bungling basic auth at this site, but firefox was ok)
```
export IAM_PASSWORD=<yourpass>
export IAM_USERNAME=<youruser>
wget --http-user=$IAM_USERNAME --http-password=$IAM_PASSWORD http://www.fki.inf.unibe.ch/DBs/iamDB/data/forms/formsA-D.tgz
wget --http-user=$IAM_USERNAME --http-password=$IAM_PASSWORD http://www.fki.inf.unibe.ch/DBs/iamDB/data/forms/formsE-H.tgz
wget --http-user=$IAM_USERNAME --http-password=$IAM_PASSWORD http://www.fki.inf.unibe.ch/DBs/iamDB/data/forms/formsI-Z.tgz
tar -zxvf formsA-D.tgz
tar -zxvf formsE-H.tgz
tar -zxvf formsI-Z.tgz
```

Next we want to turn the PNG into PDF. It's possible to go straight from PNG to PDF, but I want the JPG available too, so we take a detour.

```zsh
$ time mogrify -format jpg *.png
mogrify -format jpg *.png  868.98s user 349.32s system 291% cpu 6:58.29 total
```

The JPG size is about 30% of the PNG size, presumably due to inbuilt default quality settings or change in compression.

Now we convert to PDF. Note that a default ImageMagick installation, this will churn away then throw a policy exception at the end. The policy of no read/write PDF was put in palce to protect against an earlier sandoxing issue, but [-dSAFER violations are now fixed in gs 9.24](https://www.kb.cert.org/vuls/id/332928/). So if you are running the latest ```gs```, it is presumably back to normal and ok to mod security policies to permit PDF read/write again... So make an edit in your ImageMagick config, in the ```<policymap>``` section of ```/etc/ImageMagick/policy.xml``` (at least, for linux):

```xml
<policy domain="coder" rights="read/write" pattern="PDF" />
```

The conversion to JPG to PDF takes nearly twice as long as from PNG to JPG, coming in just under one second per page:

```zsh
$ time mogrify -format pdf *.jpg      
mogrify -format pdf *.jpg  1457.73s user 319.50s system 268% cpu 11:01.16 total
```

If the PDF files are empty, check the policy config! The PDF size should be about equal to the JPG size. 

Let's clean up with [Fred's Magick TextCleaner](http://www.fmwconcepts.com/imagemagick/textcleaner/index.php)

Put the script on your path, and enable execute permissions

The contents of the script to apply textcleaner in the whole directory are relatively simple
```
#!/bin/bash
FILES=./*.jpg
for f in $FILES
do
   g=$(echo $f | sed 's/.jpg/-clean.jpg/')
   echo "Clean $f to $g"
   textcleaner $f $g
done
```

