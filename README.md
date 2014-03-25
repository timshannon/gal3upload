# An uploader for Gallery3 written in Go.  #

There doesn't seem to be any functioning linux client uploaders and the built in Flash uploader doesn't work well so I created this to upload images from the command line. 


## Usage ##

### Quick Start ###

`./gal3upload -u "http://url-of-your-gallery/index.php" -a "apikey" -r`

This will recursively run against the current directory creating galleries for every subfolder found (if a matching gallery name doesn't already exist), and uploading any images it finds (if a matching image doesn't already exist in the gallery).

This should allow you to simply manage your gallery online by managing a local folder structure.  If you want more control see the next section for a breakdown of all options.

The default is to upload the current working directory, so there is no argument to specify the upload action. 

# Arguments #

| Argument Name | Purpose |
| ------------- | ------- |
| -u | URL of the gallery, must have http at the beginning and end in index.php | 
| -a | API key.  Look this up in the Users/Groups page of your gallery dashboard | 
| -l | List Gallery Contents. If no parent is specified with -p, then the root gallery (id 1) is used.  Albums are listed like: `Album Name [ID]` with the heirarchy represented by indentation | 
| -v | Verbose.  Adds in complete entity detail for the gallery listing |
| -p | Set parent album name to use in either the the upload, or the gallery listing (-l). Name matching is case sensitive.  Note using an the album id with -pid is faster as the name doesn't need to be looked up.|
| -pid | Set parent album id (REST ID) for use with either the upload or the gallery listing (-l) |
| -r | Recurse folder if uploading, automatically creating galleries of the subfolders.  If -l(listing) a gallery, then it will recurse the gallery displaying all sub galleries and contents. Take care using the recurse option as it may take a long time with large galleries.|
| -c | Create gallery with the name passed into the argument.  Can be used with the -p argument. |
| -f | Create a local folder structure to match the album structure in the gallery.  Usually run with -r to recursively build the folder structure.  Can be run with -p to specify the parent to build from.  If no parent is specified, it'll start at the root. |
| -rebuild | Forces a rebuild of the local cache files that stores information on the gallery structure |
| -wd | Sets the working directory.  By default it uses the directory it's executed from.  You can set this to a different directory using -wd. Note that the working directory also applies to cache files, so if you change your working directory, you may end up waiting for the local cache file to rebuild.|
| -t | Number of threads to use.  Defaults to 1.  Increasing threads will speed up uploads, but setting the threads too high can tax your web server, or your client. |
| -skipCache   | Skips building the local cache file of the gallery structure. |
| -connectFile | file path to a file containing two lines.  The first is the url, and the second is the api key.  |

## Local Cache Files ##

gal3upload will try to utilize a local cache file as much as possible.  This will speed up lookups greatly, and allow a user to work with galleries by the cached name (as the REST api only knows galleries REST urls).  gal3upload will first try to lookup a name in the local cache, and if it's not found, it will rebuild the gallery cache file from scratch.  For large galleries you will see slow performance initially, but things should speed up greatly after that.  You can force a rebuild of the local cache file by calling the -rebuild argument.  But keep in mind that any typos when working with an album name could result in extended waiting due to that name not being found in the local cache.

For uploads a local .uploadcache is created containing information on every file upload from that directory.  This means you can run the same upload script against the same tree of folders as often as you like, and you won't upload duplicate images.  The idea being that you simply have a local copy of what your online gallery looks like, and new files will automatically be detected and uploaded as you copy them into the folders.

## Examples ##

*List the contents of the gallery named "Test Gallery"*
 `./gal3upload -u "http://example.com/gallery3/index.php" -a "1234" -l -p "Test Gallery"`

*Upload the contents of the current folder to the gallery with ID 123*
 `/path/to/executable/gal3upload -u "http://example.com/gallery3/index.php" -a "1234" -pid 123`
