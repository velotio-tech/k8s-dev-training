# Assignment 4

QEMU is machine emulator and virtualizer.

QCOW2 is a storage format for QEMU copy-on-write.

## Difference virtual size and disk size
virtual size - Size of the virtual disk set when creating or expanding. 
disk size - It is the current size of the disk file. i.e how much disk space on physical server is occupied(This only applies to the )

## Backing chain:
- Backing chain concept is similar to the taking a snapshot in of the current state of the virtual machine.
- So, this snapshotting mechanism is acheieved by QEMU/KVM using a backing file.
- Backing file is paired with snapshot disk image of virtual machine. And any further writes to disk image are redirected to snapshot file as redirect-on-write.

### commands:
qemu virtual-machine boot command.
```
-> qemu-system-x86_64 -smp cores=2 -m 1024 -name TinyLinux -drive file=./Fedora-Cloud-Base-35-1.2.x86_64.qcow2
prasad@empid21060:~/my_data$ qemu-img create -f qcow2 -b Fedora-Cloud-Base-35-1.2.x86_64.qcow2 snapshot-fedora.qcow2
Formatting 'snapshot-fedora.qcow2', fmt=qcow2 size=5368709120 backing_file=Fedora-Cloud-Base-35-1.2.x86_64.qcow2 cluster_size=65536 lazy_refcounts=off refcount_bits=16
```
<br>
Backing file command.
```
prasad@empid21060:~/my_data$ qemu-img info snapshot-fedora.qcow2 
image: snapshot-fedora.qcow2
file format: qcow2
virtual size: 5 GiB (5368709120 bytes)
disk size: 196 KiB
cluster_size: 65536
backing file: Fedora-Cloud-Base-35-1.2.x86_64.qcow2
Format specific information:
    compat: 1.1
    lazy refcounts: false
    refcount bits: 16
    corrupt: false

prasad@empid21060:~/my_data$ cp Fedora-Cloud-Base-35-1.2.x86_64.qcow2 backup.qcow2
prasad@empid21060:~/my_data$ qemu-img rebase -b backup.qcow2 snapshot-fedora.qcow2 
prasad@empid21060:~/my_data$ qemu-img commit snapshot-fedora.qcow2 
Image committed.
```
<br>
Virtual disk expanding and resizing command.
```
prasad@empid21060:~/my_data$ qemu-img info snapshot-fedora.qcow2 
image: snapshot-fedora.qcow2
file format: qcow2
virtual size: 3 GiB (3221225472 bytes)
disk size: 196 KiB
cluster_size: 65536
backing file: backup.qcow2
Format specific information:
    compat: 1.1
    lazy refcounts: false
    refcount bits: 16
    corrupt: false
prasad@empid21060:~/my_data$ qemu-img resize snapshot-fedora.qcow2 +5G
Image resized.
prasad@empid21060:~/my_data$ qemu-img info snapshot-fedora.qcow2 
image: snapshot-fedora.qcow2
file format: qcow2
virtual size: 8 GiB (8589934592 bytes)
disk size: 200 KiB
cluster_size: 65536
backing file: backup.qcow2
Format specific information:
    compat: 1.1
    lazy refcounts: false
    refcount bits: 16
    corrupt: false
prasad@empid21060:~/my_data$ qemu-img resize snapshot-fedora.qcow2 -5G
qemu-img: warning: Shrinking an image will delete all data beyond the shrunken image's end. Before performing such an operation, make sure there is no important data there.
qemu-img: Use the --shrink option to perform a shrink operation.
prasad@empid21060:~/my_data$ qemu-img resize --shrink snapshot-fedora.qcow2 -4G
Image resized.
prasad@empid21060:~/my_data$ qemu-img info snapshot-fedora.qcow2 
image: snapshot-fedora.qcow2
file format: qcow2
virtual size: 4 GiB (4294967296 bytes)
disk size: 196 KiB
cluster_size: 65536
backing file: backup.qcow2
Format specific information:
    compat: 1.1
    lazy refcounts: false
    refcount bits: 16
    corrupt: false
```