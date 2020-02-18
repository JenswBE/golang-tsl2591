DEVICE="pi@192.168.127.202"
echo "Building executable"
env GOOS=linux GOARCH=arm GOARM=5 go build -o tsltest cmd/tsl2591/tsl2591.go

echo "copying files"
scp tsltest ${DEVICE}:~/tsltest

echo "setting permissions"
ssh ${DEVICE} chmod u+x ./tsltest

#echo "running program"
#ssh ${DEVICE} ./tsltest
