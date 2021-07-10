#!/bin/sh

RETRIES=30
RESULT=1

while [ $RETRIES -gt 0 ] && [ $RESULT -ne 0 ]; do
   psql -h $PG_HOST -U $PG_USER -d $PG_DATABASE -p $PG_PORT -c "select 1" > /dev/null 2>&1
   if [ $? -eq 0 ]
   then
      RESULT=0
   else
      echo "Waiting for postgres server to start, $((RETRIES)) remaining attempts..."
      RETRIES=$((RETRIES-=1))
      sleep 1
   fi
done

if [ $RESULT -eq 0 ]
then
   echo 'postgresql started'
else
   echo 'failed to start postgresql server'
   exit $RESULT
fi