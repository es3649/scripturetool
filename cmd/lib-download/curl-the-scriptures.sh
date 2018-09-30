#!/bin/bash
lang=eng

# for each testament...
for tome in ot nt bofm dc-testament pgp ; do
	# if we haven't downloaded the list of books
	if [[ ! -e lib/$tome.html ]] ; then
		# download them
		link=https://www.lds.org/scriptures/$tome?lang=$lang
		echo "downloading $link"
		curl -s $link > lib/$tome.html
	fi
	# for each book the page lists
	for book in $(grep -oP "href=\"https://www.lds.org/scriptures/$tome/[-/a-z0-9]+?\?lang=$lang\"" lib/$tome.html) ; do
		# cut the name of the book out
		book=$(echo $book | cut -d'/' -f6- | cut -d'?' -f1)
		# curl the book
		if [[ ! -e lib/$book ]]
		then
			mkdir -p lib/$book
			link=https://www.lds.org/scriptures/$tome/$book?lang=$lang
			echo "downloading $book: $link"
			curl -s $link > lib/$book.html
		fi
		# if there's more than one chapter in the book
		if [[ ! $(echo $book | grep /1 ) ]]
		then
			#download them all
			for chapter in $(grep -oP "https://www.lds.org/scriptures/$tome/$book/[0-9]+?\?lang=eng" lib/$book.html)
			do
				# cut out the chapter number
				chapnum=$(echo $chapter | cut -d'/' -f7 | cut -d'?' -f1)
				if [[ ! -d lib/$book/$chapnum ]]
				then
					# make a directory for the chapter and curl the chapter
					mkdir -p lib/$book/$chapnum
					link=https://www.lds.org/scriptures/$tome/$book/$chapnum?lang=$lang
					echo "downloading chapter $chapnum: $link"
					curl -s $link > lib/$book/$chapnum.html
				fi
			done
		fi
		# process the chapters
		
	done


done

#rm -r *.html
