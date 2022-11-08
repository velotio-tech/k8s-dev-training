someshbhalsing@Velotios-MacBook-Air k8s-dev-training % mkdir qemu

someshbhalsing@Velotios-MacBook-Air k8s-dev-training % cd qemu 

<!-- Create an empty qcow -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-img create -f qcow2 alpine.qcow2 8G 
Formatting 'alpine.qcow2', fmt=qcow2 cluster_size=65536 extended_l2=off compression_type=zlib size=8589934592 lazy_refcounts=off refcount_bits=16

<!-- start the qemu with alpine iso image and attach the alpine.qcow to it -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-system-x86_64 -m 512 -boot d -cdrom alpine-standard-3.16.2-x86_64.iso -hda alpine.qcow2 -display default -vga virtio -machine type=q35

<!-- boot the alpine vm -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-system-x86_64 -m 512 -hda alpine.qcow2

<!-- image info -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-img info alpine.qcow2 
image: alpine.qcow2
file format: qcow2
virtual size: 8 GiB (8589934592 bytes)
disk size: 352 MiB
cluster_size: 65536
Format specific information:
    compat: 1.1
    compression type: zlib
    lazy refcounts: false
    refcount bits: 16
    corrupt: false
    extended l2: false

<!-- create a snapshot -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-img snapshot -c snapshot-1 alpine.qcow2

<!-- list the snapshot -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-img snapshot -l alpine.qcow2   
Snapshot list:
ID        TAG               VM SIZE                DATE     VM CLOCK     ICOUNT
1         snapshot-1            0 B 2022-11-08 19:23:00 00:00:00.000          0

<!-- create an image from block -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-img create -f qcow2 -b alpine.qcow2 -F qcow2 alpine-copy.qcow2
Formatting 'alpine-copy.qcow2', fmt=qcow2 cluster_size=65536 extended_l2=off compression_type=zlib size=8589934592 backing_file=alpine.qcow2 backing_fmt=qcow2 lazy_refcounts=off refcount_bits=16

<!-- new image info -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-img info alpine-copy.qcow2 
image: alpine-copy.qcow2
file format: qcow2
virtual size: 8 GiB (8589934592 bytes)
disk size: 196 KiB
cluster_size: 65536
backing file: alpine.qcow2
backing file format: qcow2
Format specific information:
    compat: 1.1
    compression type: zlib
    lazy refcounts: false
    refcount bits: 16
    corrupt: false
    extended l2: false

<!-- apply a commit -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-img commit alpine-copy.qcow2 
Image committed.

<!-- check the diff between two images -->
someshbhalsing@Velotios-MacBook-Air qemu % qemu-img compare alpine.qcow2 alpine.qcow2 
Images are identical.

