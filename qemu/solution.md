#### QEMU

- free and open-source emulator
- emulates the machine's processor through dynamic binary translation
- provides a set of different hardware and device models for the machine
- enables machine run a variety of guest operating systems.
- better performance as compared to virtualbox

#### QCOW2

- file format used by the qemu images
- delays allocation of storage until it is actually needed

#### Virtual size vs Disk size

- virtual size is the size chosen while creating the drive.
- disk size is the size that is actully allocated

#### Backing chain

- chains and overlays are created for snapshotting the image.
- make it possible to apply any particular snapshot OR merge multiple images with same backing chain.