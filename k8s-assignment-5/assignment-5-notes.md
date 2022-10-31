### Qemu-img
- QEMU is defined as a generic and open source machine emulator and virtualizer. This means that Qemu can be used as machine emulator to run Operating systems and programs for one machine on a different machine. An example is running ARM program on x86 PC. In this case, QEMU can use other hypervisors like Xen or KVM for CPU extensions, to achieve what is commonly referred to as Hardware Assisted Virtualization

- QEMU performance is better than Virtual Box 

- qemu-img allows you to create, convert and modify images offline. It can handle all image formats supported by QEMU.
	   
- qemu-img is the command line utility that’s used to convert various file systems used by hypervisors like Xen, KVM, VMware, VirtualBox. qemu-img is used to format guest images, add additional storage devices and network storage e.t.c. 

### libguestfs
- libguestfs is a set of tools for accessing and modifying virtual machine (VM) disk images. You can use this for viewing and editing files inside guests, scripting changes to VMs, monitoring disk used/free statistics, creating guests, P2V, V2V, performing backups, cloning VMs, building VMs, formatting disks, resizing disks, and much more.

### QCOW2 file
- A QCOW2 file is a disk image saved in the second version of the QEMU Copy On Write (QCOW2) format, which is used by QEMU virtualization software. It stores the hard drive contents of a QEMU virtual machine. 
- You can mount a QCOW2 disk image and use it to create a virtual machine in QEMU (cross-platform) and some other virtualization programs.

### Backing chain
- Backing chain is like taking a snapshot in of the current state of the virtual machine.
- Backing files and overlays are extremely useful to rapidly instantiate thin-privisoned virtual machines.
- It is quite useful in development & test environments, so that one could quickly revert to a known state & discard the overlay.

### QEMU commands:
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img create -f raw ubuntu.img 10G
Formatting 'ubuntu.img', fmt=raw size=10737418240
```

```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img info ubuntu.img 
image: ubuntu.img
file format: raw
virtual size: 10 GiB (10737418240 bytes)
disk size: 10 GiB
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img create -f vmdk ~/ubuntu.vmdk 10G
Formatting '/home/ashish/ubuntu.vmdk', fmt=vmdk size=10737418240 compat6=off hwversion=undefined
```

```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img info ubuntu.vmdk
image: ubuntu.vmdk
file format: vmdk
virtual size: 10 GiB (10737418240 bytes)
disk size: 1.32 MiB
cluster_size: 65536
Format specific information:
    cid: 1381005449
    parent cid: 4294967295
    create type: monolithicSparse
    extents:
        [0]:
            virtual size: 10737418240
            filename: ubuntu.vmdk
            cluster size: 65536
            format: 
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img convert -O qcow2 ~/ubuntu.vmdk ~/ubuntu.qcow2
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ ls -lah ubuntu.qcow2 
-rw-r--r-- 1 ashish ashish 193K Oct 20 12:54 ubuntu.qcow2
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img info ubuntu.qcow2 
image: ubuntu.qcow2
file format: qcow2
virtual size: 10 GiB (10737418240 bytes)
disk size: 204 KiB
cluster_size: 65536
Format specific information:
    compat: 1.1
    lazy refcounts: false
    refcount bits: 16
    corrupt: false
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img check ubuntu.vmdk 
No errors were found on the image.
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img check ubuntu.img
qemu-img: This image format does not support checks
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img check ubuntu.qcow2 
No errors were found on the image.
Image end offset: 262144
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img snapshot -c ubuntu_fresh ubuntu.qcow2
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img snapshot -l ubuntu.qcow2 
Snapshot list:
ID        TAG                     VM SIZE                DATE       VM CLOCK
1         ubuntu_fresh                0 B 2022-10-20 12:56:34   00:00:00.000
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img snapshot -d 1 ubuntu.qcow2 
qemu-img: Could not delete snapshot '1': snapshot not found
```
```
ashish@velotio-ThinkPad-E14-Gen-2:~$ qemu-img snapshot -d ubuntu_fresh ubuntu.qcow2 
```