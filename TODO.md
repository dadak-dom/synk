# TODO list: Bugs, new features, notes, etc.

## Known Bugs:

### <span style="color: red">Copying links-to-directories</span>

Copying a file that is a "Link to Folder" results in:

- Nothing when that file is directly within the shared folder
- If the file is in a subdirectory of the shared folder, the program tries to copy it over, then realizes it's a directory, thus crashing the program

### <span style="color: orange">Peer list gets forgetten when changing pages</span>

Per the title. Find a way to have the peer list persist. Also have it remember the selected peers

### <span style="color: orange">Need dynamic folder/file updates in shraed folder view</span>

Per the title, when a list of folders/files is listed, any changes to the contents should be reflected without needing to change pages

## New features to add:

Sorted and color-coded by priority

### <span style="color: red">TOP: Approved network list</span>

On startup, check the network name and compare it to the approved network list

### <span style="color: red">TOP: Notify user when transfer done!</span>

Per the title.

### <span style="color: orange">(1) Choose files to ignore</span>

In the folder selector menu, list the files in the currently chosen directory. For each file, if they are clicked, toggle a strike-through on their name. For every file/subfolder that is selected (e.g. <s>test.txt</s> or <s>folder</s>), exempt it from the synk process.<br>
Almost done, just need to test and make sure that it works + maybe add subfolder support?<br>
Make the changes to the list propogate to peers<br>

### <span style="color: orange">(2) New config option: Launch on start</span>

In settings menu, allow a user to choose if the app should launch on startup. <br>
Will need new config options on the backend and set up communication between FE and BE

### <span style="color: orange">(3) New config option: Auto-synk</span>

In settings menu, user should have option to auto-synk. If selected, this mode should do the following:

1. Set up a watcher to watch the shared directory. <br> Make sure if the directory changes, the watcher changes accordingly
2. If a change is logged, check if there are peers on the network. Verify that peers in the peerList are still around, maybe ping them?
3. If they (peers) are around, initiate a "synk" without the user's input
4. Continue as long as the app is running

### <span style="color: orange">(4) Make app dockable</span>

Per the title, see if I can make the app dockable (e.g. like Discord)

### <span style="color: cornflowerblue">(5) Add some security layer</span>

In the interest of user safety, I should probably add some sort of encryption / safety mechanism to the requests being sent. This is not a top priority because the idea is that the user would use this on their home network, and so they wouldn't be susceptible to "man in the middle" attacks or anything of that nature. However, it would be nice to give the user some peace-of-mind in case they waned to synk their files when connected to a network they may not have full control over, e.g. in the office or at a coffee shop.

### <span style="color: cornflowerblue">(6) Create mobile app version</span>

I suppose this would mean making sure that the CSS looks good on mobile devices. Mobile builds are not [implemented yet](https://github.com/wailsapp/wails/issues/1481), but mi.ght be in the near future

LATER:

- Make the navbar pretty.
- Overall, think of a cohesive aesthetic to create for the app, and then make the assets for it
- Expanding on (\*): show a comparison of the files on the local vs. on the remote machine before doing the synk (as a sort of confirmation).
  Would need to create a graphical rendering of the filesystem (nested folders, etc.)

BACKEND:
NOW:

- Generate a random string of words or numbers on startup. Use the FRONTEND to confirm that you want to synk with the computer with that ID
- only launch the API when actually running a synk. (Less drastic measure: toggle the API to off when not synking.)

LATER (LONG TERM):

- Find a way to allow users to automatically synk so long as both computers are on
- Find a way to "remember" other devices and automatically synk with them

- Folder Selector:
  - make it prettier
- Multicast
  - Make sure that the library works when used on the same LAN network
- File sharing
  - code up the file sending / receiving logic

FIXME:

- Make it so that the shared folder is the default when opening the folder selector
- For the time being, folders don't work. Will need to implement some sort of check that creates a subdirectory if it doesn't already exist
- There is a bug (on windows) where the selected local IP for the http server is on a different subnet. Will need to make it so that the API is only launched once a connection has been made via multicast
- Will need to make it so that the program is always listening for new peers in the background. Will have to rejigger the logic a little bit

NOTES:

- For some reason, the multicasting only works between windows/linux when I change the multicast address (e.g. from 224.0.0.0 to 224.0.0.1)

NOTE!:

- On Linux builds, DO NOT use CSS rules like this:
  .outer-rule {
  .inner-rule
  }

This breaks the inner rule
