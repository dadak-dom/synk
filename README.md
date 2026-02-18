# Synk 🔄

A peer-to-peer, LAN file sharing program 📂. Built for seamless transitions between your working devices.

## About ❔

<b>Synk</b> is a desktop app designed for sharing files across a private network. <b>Synk</b> doesn't just allow you to send files; it automatically "synk"-chronizes the contents of your selected <b>Shared Folder</b>. For example:<br><br> On my <b>Desktop</b> 🖥️, say you have:

```golang
your_shared_folder/
  homework.txt
  resume.pdf
  pictures/
    img001.png
    img002.png
```

<br> On your <b>Laptop</b> 💻, you added some notes to your homework, and maybe found a good study guide:

```golang
your_shared_folder/
  homework.txt // <- newer than the version on your Desktop
  study_guide.pdf
  resume.pdf
  pictures/ // <- maybe you took some more pictures, too.
    img003.png
    img004.png
```

After opening the application on both devices, selecting a peer, and clicking the big Synk (🔄) button, <b>both</b> devices should now contain:

```golang
your_shared_folder/
  homework.txt // updated to newest version!
  study_guide.pdf
  resume.pdf
  pictures/
    img001.png
    img002.png
    img003.png
    img004.png
```

The idea is for <b>Synk</b> to allow you to seamlessly transition your work from one device to another. With this program, you can quickly and easily send files to yourself without relying on third-party services or setting up a physical connection.

## Development 👩‍💻

<b>Synk</b> is built with <a href="https://wails.io/">Wails</a>. Follow their instructions for setting it up.

Run the app locally with `wails dev`<br>
Compile the app with `wails build`
