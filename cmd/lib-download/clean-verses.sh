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
		heading=$(grep -P '\<p class=\"study-summary\".+?\</p\>' $chapter | perl -pe "s@\<[^>]+\>@@g; s@\t+@@g" | tr "\n" "\\n") # This last pipe might not actually work...

		# add it to the raw text files
		# echo $heading >> $book$chapnum.text

		# get all the verses from the html
		# returns html lines with '~' instead of ' ' 
		VERSES=$(grep -oP "\<p class=\"verse\".+?\</p\>" $chapter | perl -pe 's/ /~/g')

		# print out chapter info in json
		echo "{" >> $book$chapnum.json
		echo "    \"book\":\"$bookname\"," >> $book$chapnum.json
		echo "    \"chapter\":\"$chapnum\"," >> $book$chapnum.json
		echo "    \"heading\":\"$heading\"," >> $book$chapnum.json
		echo "    \"verses\":" >> $book$chapnum.json
		echo "    {" >> $book$chapnum.json

		# get some verse metrics for comma s"pacing
		VS_COUNT=$(echo $VERSES | tr " " "\n" | wc -l)
		VS_Counter=0
		# echo "Parsing through $VS_COUNT verse(s)"

		# for each verse:
		for verse in $VERSES ; do
			# clean out the transcript
			text=$(echo "$verse" | perl -pe "s@>(.)</sup>@>{\1}@g; s@<[^>]*>@@g; s@^([0-9]{1,3}) @\1@g")

			# print the number (field 1), and take the number off the text
			versenum=$(echo $text | cut -d'~' -f1)
			text=$(echo $text | cut -d'~' -f2-)

			# this awk script needs to do a ton of stuff...
			# find all the footnotes, it DOES NOT need to capture them
			# it needs to count the number of charactes from the beginning of the line until the footnote NOT including the note markers
			# it needs to return these numbers in some form, a list maybe?
			# the first arg of the list is the count of other args in the list
			positions=$(echo $text | awk --field-separator="{.}" '
				{ res = sprintf(NF);
					CT=0; 
					for (i=1; i < NF; i = i+1) {
						CT=CT+length($i);
						res=res" "sprintf(CT);
					}
				} 
				END {print res}' | tr "\n" " ")
			# echo $positions

			# clean out the refs and put spaces back before we print the text to the file
			text=$(echo $text | perl -pe 's/~/ /g; s/{.}//g')
			echo "        \"$versenum\":" >> $book$chapnum.json
			echo "        {" >> $book$chapnum.json	
			echo "            \"text\":\"$text\"," >> $book$chapnum.json
			echo "            \"footnotes\":" >> $book$chapnum.json
			echo "            [" >> $book$chapnum.json
			
			# clean out the footnotes
			# get a list of the notes
			REF_LIST=$( echo $verse | grep -oP "https://www.lds.org/scriptures/footnote\?lang=eng&amp;data-uri=/scriptures/.+?/.+?/[0-9]+&amp;noteID=note[0-9]+." )
			REF_LIST=$( echo "$REF_LIST" | perl -pe "s@&amp;@&@g" )

			# download them all
			# no worry about indicies for zero cases. Zero cases are handled, because this loop doesn't run if there are no refs in $REF_LIST
			# echo "Downloading references for $bookname $chapnum:$versenum"
			pos=2
			for ref in $REF_LIST ; do

				PAGE=$(curl -s $ref)
				#clean them
				NOTE=$( echo "$PAGE" | tr "\t" " " | perl -pe "s@\<p id=\"note[0-9]+(.)@[\1] <@g; s@\<[^>]+?\>@@g" | tr "\n" " " | perl -pe "s@ {2,}@ @g")
				#NOTE=$( echo $PAGE | perl -pe "s@\<p id=\"note[0-9]+(.)@[\1] <@g" )
				#echo "$NOTE"
				#NOTE=$( echo $NOTE | perl -pe "s@\<[^>]+?\>@@g" )
				#echo "$NOTE"
				#NOTE=$( echo $NOTE | perl -pe "s@ {2,}@ @g" )
				#echo "$NOTE"
				echo "                {" >> $book$chapnum.json
				NOTE_POSITION=$(echo $positions | cut -d" " -f$pos)
				echo "                    \"position\":$NOTE_POSITION," >> $book$chapnum.json
				echo "                    \"reference\":\"$NOTE\"" >> $book$chapnum.json

				# only put the comma in if this is not the last ref
				# $pos[1] holds the number of refs
				if [[ ! "$pos" == "$(echo $positions | cut -d" " -f1)" ]]; then
					echo "                }," >> $book$chapnum.json
				else
					echo "                }" >> $book$chapnum.json
				fi
				
				let pos++
			done


			echo "            ]" >> $book$chapnum.json

			# update the counter, then insert the corresponding item ending.
			let VS_Counter++
			if [[ "$VS_Counter" == "$VS_COUNT" ]]; then
				echo "        }" >> $book$chapnum.json
			else
				echo "        }," >> $book$chapnum.json
			fi
		done

		# put the end of the json in place
		echo "    }" >> $book$chapnum.json
		echo "}" >> $book$chapnum.json
		
	done
done
