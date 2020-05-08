productiondir=/var/www
productionname=samuellando.com
builddir=build
productionip=samuellando.com

echo "Creating build tarball"
tar cf build.tar build/*
echo "Removing old production build"

ssh -o StrictHostKeyChecking=no deployer@$productionip "rm -r $productiondir/$productionname"
echo "Copying new build tarball to production"
scp -o StrictHostKeyChecking=no build.tar deployer@$productionip:$productiondir/build.tar
echo "Unpacking new production build"
ssh -o StrictHostKeyChecking=no deployer@$productionip "tar xf $productiondir/build.tar -C $productiondir; mv $productiondir/build $productiondir/$productionname; rm $productiondir/build.tar"
