# Stor

Stor is a very simple store/restore, in-memory caching utility, with a very simple TCP socket interface.

It is a go-experimentation project, and thus really simple and perhaps immature.

It is currenly being refactored into a more senseful architecture in the "next" branch.

## Run up the server

Read file "HACKIG" for building instructions and how to start up the server. The server will run in the foreground and can be stopped by sending it the INT signal.

## Stor and restor

Stor supports five basic commands with very simple syntax:

* **set** _key_ **\r\n** _value_ **\r\n**: store value _value_ under key _key_. Notice there is a maximum size for keys and values (in the parser code). Server response: **STORED**.
* **get** _key,_ _..._ **\r\n**: retrieve values under each _key_, previously stored with a _set_ command. Server response **VALUE** _key_ **\r\n** _value_ for each _key_.
* **delete** _key_ **\r\n**: delete the entry under _key_. Server response: **DELETED** or **NOT_FOUND** followed by **\r\n**.
* **stats** **\r\n**: request some server stats, like number of get-s, set-s, cache hits and misses, etc. Server response: all these followed by **END** **\r\n**.
* **quit** **\r\n**: signal end-of-communication with the server. The server will close the connection after this. Server response: none.

**NOTICE** that in order to read the server's responses without I/O blocks, one should always send the _quit_ command first.

## Example in Ruby

Fire up irb

    $ irb -r socket
    irb(main):001:0> def stor; TCPSocket.new 'localhost', 11212 end
    => :stor
    irb(main):002:0> s = stor
    => #<TCPSocket:fd 9>
    irb(main):003:0> s.write "set one \r\n This is text one, yay \r\n"
    => 35
    irb(main):004:0> s.write "set two \r\n This is text two \n new line, hooray! \r\n"
    => 50
    irb(main):005:0> s.write "set three \r\n I need to see \r\n"
    => 29
    irb(main):006:0> s.write "get one two three four \r\n"
    => 25
    irb(main):007:0> for key in %w[ one three four] do s.write "delete #{key} \r\n" end
    => ["one", "three", "four"]
    irb(main):008:0> s.write "get one two three four \r\n"
    => 25
    irb(main):009:0>  s.write "stats \r\n"
    => 8
    irb(main):010:0> s.write "quit \r\n"
    => 7
    irb(main):011:0> puts s.read
    STORED
    STORED
    STORED
    VALUE one
    This is text one, yay
    VALUE two
    This is text two
     new line, hooray!
    VALUE three
    I need to see
    END
    DELETED
    DELETED
    NOT_FOUND
    VALUE two
    This is text two
     new line, hooray!
    END
    cmd_get 8
    cmd_set 3
    get_hits 4
    get_misses 4
    delete_hits 2
    delete_misses 1
    curr_items 1
    limit_items 65535
    END
    => nil
    irb(main):012:0> s.close
    => nil
    irb(main):013:0>
