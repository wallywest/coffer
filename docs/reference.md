# Reference

Assets are currently stored in MongoDB GridFS.  The GridFS collections for the assets default
to the name *vcsfs.files* and *vcs.chunks*


### Recording metatdata 

### Asset Document Schema

Here is a what a typical record will look like in GridFS for a recording metadata.

```json
{
  "_id" : ObjectId("56b8e53d31b2f50001744a40"),
  "length" : NumberLong(29964),
  "chunkSize" : 261120,
  "uploadDate" : ISODate("2016-02-08T18:58:05Z"),
  "filename" : "13d84d3d993e788929a004292a80afd2df434605",
  "metadata" : {
    "accountId" : "AC56445f9d0b977d270d02b7026719484c2b6bf369",
    "callId" : "CA465c118caa4ed1ba2b91c62e6a4d3f9ee8d3cb8c",
    "duration" : 8,
    "fileId" : "13d84d3d993e788929a004292a80afd2df434605",
    "fileName" : "callrec_0_O287_172.20.152.36_5237_1.0.14_1454604276.wav",
    "fileSize" : 29964,
    "mimeType" : "audio/wav",
    "downloadUrl" : "/Accounts/AC56445f9d0b977d270d02b7026719484c2b6bf369/recordings/13d84d3d993e788929a004292a80afd2df434605"
  }
}
```


|Key | Description |
| --- | --- |
| _id | ObjectId given to asset by MongoDB |
| length | length of the asset file |
| uploadDate | date the file was stored |
| fileName | unique id generated for asset, same as fileId |
| metadata.accountId | accountId associated with the asset |
| metadata.callId | callId associated with the asset |
| metadata.duration | duration of the asset |
| metadata.fileId | unique id generated for asset, same as fileName |
| metatdata.fileSize | size of the recorded asset, same as length |
| metadata.mimeType | mimetype of the asset.  Only supports audio/wav currently |
| metatdata.downloadUrl | internal address mapping for an asset.  Used by mediaserver |


Coffer APIs
===========

These endpoints are assumed to be forwarded internally from the rest public api


### Get a Recording: ``GET /Account/:accountId/Recordings/:recordingId``

### Download and Asset: ``GET /Account/:accountId/Recordings/:recordingId/Download``

### List Recordings: ``GET /Account/:accountId/Recordings``

### Delete a Recording: ``DELETE /Account/:accountId/Recordings``
