command rm bin/api
command echo "Creating image with version $1"
command make PLATFORM=linux/amd64
command rsync -rav $(pwd)/bin/* movies-api:/home/ec2-user
command echo "Files were deployed"