#!/bin/bash

# desired text file format:

#
#book heading this is the heading for this chapter of this book. stuff happens here
#book chap:1 Lorem ipsum [a]dolor sit [b]amet... 
#<a "book chap:vs book chap:vs"> <b "GR an insult">
#book chap:2 ....etc.


for book in $(ls -d lib/*/) ; do
	for chapter in $(ls $book*.html) ; do
		if [[ ! -e $chapter ]] ; then
			continue
		fi
		# regex the chapter number out
		bookname=$(echo $book | cut -d'/' -f2)
		chapnum=$(echo "$chapter" | cut -d'.' -f1 | cut -d'/' -f3)

		# remove the old archive files
		# rm $book$chapnum.text
		# rm $book$chapnum.note
		echo "cleaning $book$chapnum"

		# get the chapter heading and format it
		heading=$(grep -P '\<p class=\"study-summary\".+?\</p\>' $chapter | perl -pe "s@\<[^>]+\>@@g; s@\t+@@g")

		# add it to the raw text files
		# echo $heading >> $book$chapnum.text

		# get all the verses from the html
		# returns html lines with '~' instead of ' ' 
		VERSES=$(grep -oP "\<p class=\"verse\".+?\</p\>" $chapter | perl -pe 's/ /~/g')

		# print out chapter info in json
		echo "{" >> $chapnum.json
		echo "    \"book\":\"$bookname\"," >> $chapnum.json
		echo "    \"chapter\":\"$chapnum\"," >> $chapnum.json
		echo "    \"heading\":\"$heading\"," >> $chapnum.json
		echo "    \"verses\":" >> $chapnum.json
		echo "    {" >> $chapnum.json

		# for each verse:
		for verse in $VERSES ; do
			# clean out the transcript
			text=$(echo "$verse" | perl -pe "s@>(.)</sup>@>{\1}@g; s@<[^>]*>@@g; s@^([0-9]{1,3}) @\1@g")

			# print the number (field 1)
			echo "        \"$(echo $text | cut -d' ' -f1)\":" >> $chapnum.json
			echo "        {" >> $chapnum.json
			echo "            \"text\":$\"(echo $text | cut -d' ' -f2-)\"," >> $chapnum.json
			echo "            \"footnotes\":" >> $chapnum.json
			echo "            [" >> $chapnum.json

			# this awk script needs to do a ton of stuff...
			# find all the footnotes, it DOES NOT need to capture them
			# it needs to count the number of charactes from the beginning of the line until the footnote NOT including the note markers
			# it needs to return these numbers in some form, a list maybe?
			positions=$(echo $text | awk --field-separator="{.}" '
				BEGIN { res = "" }
				{ CT=0; 
					for (i=1; i < NF; i = i+1) {
						print $i "\n";
						CT=CT+length($i)
					}
					res=res sprintf(CT);
				} 
				END {print res}' )
			echo $positions
			
			# clean out the footnotes
			# get a list of the notes
			REF_LIST=$( echo $verse | grep -oP "https://www.lds.org/scriptures/footnote\?lang=eng&amp;data-uri=/scriptures/.+?/.+?/[0-9]+&amp;noteID=note[0-9]+." )
			REF_LIST=$( echo "$REF_LIST" | perl -pe "s@&amp;@&@g" )
			FOOTNOTES=""
			#download them all
			for ref in $REF_LIST ; do

				echo "Downloading references for $bookname Chapter $chapnum"
				PAGE=$(curl -s $ref)
				#clean them
				NOTE=$( echo "$PAGE" | tr "\t" " " | perl -pe "s@\<p id=\"note[0-9]+(.)@[\1] <@g; s@\<[^>]+?\>@@g" | tr "\n" " " | perl -pe "s@ {2,}@ @g")
				#NOTE=$( echo $PAGE | perl -pe "s@\<p id=\"note[0-9]+(.)@[\1] <@g" )
				#echo "$NOTE"
				#NOTE=$( echo $NOTE | perl -pe "s@\<[^>]+?\>@@g" )
				#echo "$NOTE"
				#NOTE=$( echo $NOTE | perl -pe "s@ {2,}@ @g" )
				#echo "$NOTE"
				FOOTNOTES="$FOOTNOTES $NOTE"
			done


			echo "            ]" >> $chapnum.json
			echo "        }," >> $chapnum.json # do something about this comma. It doesn't always go in
		done

		# put the end of the json in place
		echo "    }" >> $chapnum.json
		echo "}" >> $chapnum.json
		
	done
done
