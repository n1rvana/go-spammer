# go-spammer

Several ways of reading data off of channels are shown in this code.  Most of them block.  So, if you have a go routine managing multiple channels, you want to use select.
