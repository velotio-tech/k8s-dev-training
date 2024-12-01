# assignment 2
to run 
```make install```
and then ```make run```


This code includes a custom API added a cronjob struct and api for it.

for cronjob status added another struct CronJobStatus with declaring it as subresource

tried adding and removing the fields of the CRD and changing validations with markers

when we declare subresource for our CRD i.e scale or status we get another endpoint added in api server
like ...<crd>/status, ...<crd>/scale
PUT request to endpoint ...<crd>/status takes a custom resource object and ignore changes to anything except the status and it only validates the status stanza

PUT/POST/PATCH requests to ...<crd>/ ignores changes to the status stanza

