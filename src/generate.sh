counter=0
while [ $counter -le 99 ]
do
node ./src/index.js $counter
((counter++))
done