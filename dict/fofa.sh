for i in `cat result.txt`
do
	curl "https://fofa.so/api/v1/search/all?email=${FOFA_EMAIL}&key=${FOFA_KEY}&qbase64=${i}&size=10000&page=0" |jq '.results[][0]'|sed 's/"//g' >> query.txt
done
