**Kubernetes controllers are the *brain* (orchestration), but utilities like `qemu-img` and `libguestfs` are the *hands* (disk engine).**

When you build a K8s controller to back up or move virtual machine disks (a **Data Mover** operator), Go code alone cannot manipulate raw `.qcow2` bytes efficiently. Instead, your Go controller runs or invokes these low-level system commands under the hood.

Here is your step-by-step guide to mastering and completing Assignment 5.

---

# Assignment 5: `qemu-img` & Data Mover Utilities

---

## 1. The Core "Why": Why `qemu-img` when we have Go Controllers?

* **Go Controllers (The Brain):** Watch Kubernetes resources (CRDs like `VirtualMachineBackup`), manage status, decide *when* to trigger a backup, and orchestrate pods.
* **`qemu-img` / `libguestfs` (The Engine):** High-performance C/C++ utilities that perform low-level disk manipulation—converting disk formats, managing snapshot backing chains, and reading filesystems inside VM images without booting the VM.

Your Go controller literally executes `exec.Command("qemu-img", ...)` or uses Cgo bindings behind the scenes to perform data movement!

---

## 2. Virtual Size vs. Disk Size

* **Virtual Size (Provisioned Capacity):** The size of the disk as seen **inside the guest OS** (e.g., a "100 GB hard drive").
* **Disk Size (Physical Allocation):** The actual space consumed on the **host storage drive** right now (e.g., 2.5 GB of actual written data).

### Hands-on Exercise:

```bash
# 1. Create a 10GB qcow2 image
qemu-img create -f qcow2 test-disk.qcow2 10G

# 2. Inspect size
qemu-img info test-disk.qcow2

```

**What you will observe:**

```text
virtual size: 10 GiB (10737418240 bytes)
disk size: 196 KiB

```

The file only takes **196 KiB** on disk because `qcow2` dynamically grows as data gets written!

---

## 3. Creating & Restoring `qcow2` Images in a Data Mover

A **Data Mover** backs up block devices or filesystems into thin `.qcow2` files and restores them when needed.

```
                   ┌────────────────────────┐
                   │  Raw Block Device /    │
                   │  PV (/dev/sdb or file) │
                   └───────────┬────────────┘
                               │
            Backup (qemu-img convert)
                               │
                               ▼
                   ┌────────────────────────┐
                   │ Compressed qcow2 Image │
                   │  (backup-v1.qcow2)     │
                   └───────────┬────────────┘
                               │
           Restore (qemu-img convert)
                               │
                               ▼
                   ┌────────────────────────┐
                   │ Restored Block Device  │
                   │    or Filesystem       │
                   └────────────────────────┘

```

### A. Create a `qcow2` Image from a Block Device or File

```bash
# Convert a raw block device (e.g., /dev/sdb) or raw file to a thin, compressed qcow2 image
qemu-img convert -p -c -O qcow2 /dev/sdb backup-v1.qcow2

```

* `-p`: Shows progress percentage.
* `-c`: Enables compression (saves huge storage during backup).
* `-O qcow2`: Target output format.

### B. Restore Data from a `qcow2` Image back to a Block Device

```bash
# Restore backup image to a target block device or raw file
qemu-img convert -p -O raw backup-v1.qcow2 /dev/sdc

```

---

## 4. Backing Chains & Differential Backups (`-F` and `-F` / `-b`)

### What is a Backing Chain?

A backing chain links a **read-only base image** (e.g., Full Backup) to an **overlay/diff image** (e.g., Incremental Backup). When reading data, the system checks the overlay first; if missing, it falls back to the base image.

```
   [ Base Image: full-backup.qcow2 ] (Full Backup - Read Only)
                  ▲
                  │  (Backing File Reference)
   [ Overlay Image: diff-v1.qcow2 ]  (Incremental Changes Only)

```

### Step-by-Step Backing Chain Commands

#### 1. Create Base Image (Full Backup)

```bash
qemu-img create -f qcow2 base.qcow2 5G

```

#### 2. Create an Overlay Image (Differential/Incremental)

```bash
qemu-img create -f qcow2 -b base.qcow2 -F qcow2 diff1.qcow2

```

* `-b base.qcow2`: Specifies the parent backing file.
* `-F qcow2`: Explicitly declares the format of the backing file (security best practice).

#### 3. Creating a Differential Image using `-F` and rebase (`-b`)

If you have two disk states (Data Source A and Data Source B) and want to create a lightweight delta file referencing an existing base:

```bash
# Create diff image explicitly referencing base.qcow2
qemu-img create -f qcow2 -b base.qcow2 -F qcow2 delta.qcow2

# Change or explicitly rebind the backing chain using rebase
qemu-img rebase -u -b new-base.qcow2 -F qcow2 delta.qcow2

```

* `-u` (unsafe mode): Updates the backing file reference without copying bytes. Useful when moving or renaming backing files in storage.

---

## 5. `libguestfs` Commands in Data Movers

`libguestfs` allows you to open, modify, and inspect VM disk contents **without booting a virtual machine**.

### Essential Commands:

```bash
# 1. List filesystems inside a qcow2 backup image
virt-filesystems -a backup-v1.qcow2 --extra

# 2. Inspect filesystem disk space usage inside the guest image
virt-df -h -a backup-v1.qcow2

# 3. Mount and view files inside an image directly (guestfish)
guestfish --ro -a backup-v1.qcow2 -m /dev/sda1 cat /etc/os-release

```

---

## 6. Complete Practical Flow (Cheat Sheet)

```bash
# --- STEP 1: Create a raw test file representing a storage volume ---
dd if=/dev/urandom of=raw-volume.img bs=1M count=100

# --- STEP 2: Backup raw volume into a thin qcow2 image ---
qemu-img convert -p -c -O qcow2 raw-volume.img full-backup.qcow2

# --- STEP 3: Check Virtual vs Disk size ---
qemu-img info full-backup.qcow2

# --- STEP 4: Create an incremental diff image on top of full-backup ---
qemu-img create -f qcow2 -b full-backup.qcow2 -F qcow2 inc-backup-1.qcow2

# --- STEP 5: Inspect backing chain relationship ---
qemu-img info inc-backup-1.qcow2

# --- STEP 6: Restore full data back to a raw block volume ---
qemu-img convert -p -O raw full-backup.qcow2 restored-volume.img

```

---
