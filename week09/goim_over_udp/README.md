- 总结几种 socket 粘包的解包方式: 
  - fix length
  - delimiter based
  - length field based frame decoder。

https://netty.io/4.0/api/io/netty/handler/codec/FixedLengthFrameDecoder.html
```
A decoder that splits the received ByteBufs by the fixed number of bytes. For example, if you received the following four fragmented packets:
 +---+----+------+----+
 | A | BC | DEFG | HI |
 +---+----+------+----+
 
A FixedLengthFrameDecoder(3) will decode them into the following three packets with the fixed length:
 +-----+-----+-----+
 | ABC | DEF | GHI |
 +-----+-----+-----+
```

https://netty.io/4.0/api/io/netty/handler/codec/DelimiterBasedFrameDecoder.html
```
DelimiterBasedFrameDecoder allows you to specify more than one delimiter. If more than one delimiter is found in the buffer, it chooses the delimiter which produces the shortest frame. For example, if you have the following data in the buffer:

 +--------------+
 | ABC\nDEF\r\n |
 +--------------+
 
a DelimiterBasedFrameDecoder(Delimiters.lineDelimiter()) will choose '\n' as the first delimiter and produce two frames:
 +-----+-----+
 | ABC | DEF |
 +-----+-----+
 
rather than incorrectly choosing '\r\n' as the first delimiter:
 +----------+
 | ABC\nDEF |
 +----------+
```

https://netty.io/4.0/api/io/netty/handler/codec/LengthFieldBasedFrameDecoder.html
```
A decoder that splits the received ByteBufs dynamically by the value of the length field in the message. It is particularly useful when you decode a binary message which has an integer header field that represents the length of the message body or the whole message.
LengthFieldBasedFrameDecoder has many configuration parameters so that it can decode any message with a length field, which is often seen in proprietary client-server protocols. Here are some example that will give you the basic idea on which option does what.

2 bytes length field at offset 0, do not strip header
The value of the length field in this example is 12 (0x0C) which represents the length of "HELLO, WORLD". By default, the decoder assumes that the length field represents the number of the bytes that follows the length field. Therefore, it can be decoded with the simplistic parameter combination.

 lengthFieldOffset   = 0
 lengthFieldLength   = 2
 lengthAdjustment    = 0
 initialBytesToStrip = 0 (= do not strip header)

 BEFORE DECODE (14 bytes)         AFTER DECODE (14 bytes)
 +--------+----------------+      +--------+----------------+
 | Length | Actual Content |----->| Length | Actual Content |
 | 0x000C | "HELLO, WORLD" |      | 0x000C | "HELLO, WORLD" |
 +--------+----------------+      +--------+----------------+
 
2 bytes length field at offset 0, strip header
Because we can get the length of the content by calling ByteBuf.readableBytes(), you might want to strip the length field by specifying initialBytesToStrip. In this example, we specified 2, that is same with the length of the length field, to strip the first two bytes.

 lengthFieldOffset   = 0
 lengthFieldLength   = 2
 lengthAdjustment    = 0
 initialBytesToStrip = 2 (= the length of the Length field)

 BEFORE DECODE (14 bytes)         AFTER DECODE (12 bytes)
 +--------+----------------+      +----------------+
 | Length | Actual Content |----->| Actual Content |
 | 0x000C | "HELLO, WORLD" |      | "HELLO, WORLD" |
 +--------+----------------+      +----------------+
 
2 bytes length field at offset 0, do not strip header, the length field represents the length of the whole message
In most cases, the length field represents the length of the message body only, as shown in the previous examples. However, in some protocols, the length field represents the length of the whole message, including the message header. In such a case, we specify a non-zero lengthAdjustment. Because the length value in this example message is always greater than the body length by 2, we specify -2 as lengthAdjustment for compensation.

 lengthFieldOffset   =  0
 lengthFieldLength   =  2
 lengthAdjustment    = -2 (= the length of the Length field)
 initialBytesToStrip =  0

 BEFORE DECODE (14 bytes)         AFTER DECODE (14 bytes)
 +--------+----------------+      +--------+----------------+
 | Length | Actual Content |----->| Length | Actual Content |
 | 0x000E | "HELLO, WORLD" |      | 0x000E | "HELLO, WORLD" |
 +--------+----------------+      +--------+----------------+
 
3 bytes length field at the end of 5 bytes header, do not strip header
The following message is a simple variation of the first example. An extra header value is prepended to the message. lengthAdjustment is zero again because the decoder always takes the length of the prepended data into account during frame length calculation.

 lengthFieldOffset   = 2 (= the length of Header 1)
 lengthFieldLength   = 3
 lengthAdjustment    = 0
 initialBytesToStrip = 0

 BEFORE DECODE (17 bytes)                      AFTER DECODE (17 bytes)
 +----------+----------+----------------+      +----------+----------+----------------+
 | Header 1 |  Length  | Actual Content |----->| Header 1 |  Length  | Actual Content |
 |  0xCAFE  | 0x00000C | "HELLO, WORLD" |      |  0xCAFE  | 0x00000C | "HELLO, WORLD" |
 +----------+----------+----------------+      +----------+----------+----------------+
 
3 bytes length field at the beginning of 5 bytes header, do not strip header
This is an advanced example that shows the case where there is an extra header between the length field and the message body. You have to specify a positive lengthAdjustment so that the decoder counts the extra header into the frame length calculation.

 lengthFieldOffset   = 0
 lengthFieldLength   = 3
 lengthAdjustment    = 2 (= the length of Header 1)
 initialBytesToStrip = 0

 BEFORE DECODE (17 bytes)                      AFTER DECODE (17 bytes)
 +----------+----------+----------------+      +----------+----------+----------------+
 |  Length  | Header 1 | Actual Content |----->|  Length  | Header 1 | Actual Content |
 | 0x00000C |  0xCAFE  | "HELLO, WORLD" |      | 0x00000C |  0xCAFE  | "HELLO, WORLD" |
 +----------+----------+----------------+      +----------+----------+----------------+
 
2 bytes length field at offset 1 in the middle of 4 bytes header, strip the first header field and the length field
This is a combination of all the examples above. There are the prepended header before the length field and the extra header after the length field. The prepended header affects the lengthFieldOffset and the extra header affects the lengthAdjustment. We also specified a non-zero initialBytesToStrip to strip the length field and the prepended header from the frame. If you don't want to strip the prepended header, you could specify 0 for initialBytesToSkip.

 lengthFieldOffset   = 1 (= the length of HDR1)
 lengthFieldLength   = 2
 lengthAdjustment    = 1 (= the length of HDR2)
 initialBytesToStrip = 3 (= the length of HDR1 + LEN)

 BEFORE DECODE (16 bytes)                       AFTER DECODE (13 bytes)
 +------+--------+------+----------------+      +------+----------------+
 | HDR1 | Length | HDR2 | Actual Content |----->| HDR2 | Actual Content |
 | 0xCA | 0x000C | 0xFE | "HELLO, WORLD" |      | 0xFE | "HELLO, WORLD" |
 +------+--------+------+----------------+      +------+----------------+
 
2 bytes length field at offset 1 in the middle of 4 bytes header, strip the first header field and the length field, the length field represents the length of the whole message
Let's give another twist to the previous example. The only difference from the previous example is that the length field represents the length of the whole message instead of the message body, just like the third example. We have to count the length of HDR1 and Length into lengthAdjustment. Please note that we don't need to take the length of HDR2 into account because the length field already includes the whole header length.

 lengthFieldOffset   =  1
 lengthFieldLength   =  2
 lengthAdjustment    = -3 (= the length of HDR1 + LEN, negative)
 initialBytesToStrip =  3

 BEFORE DECODE (16 bytes)                       AFTER DECODE (13 bytes)
 +------+--------+------+----------------+      +------+----------------+
 | HDR1 | Length | HDR2 | Actual Content |----->| HDR2 | Actual Content |
 | 0xCA | 0x0010 | 0xFE | "HELLO, WORLD" |      | 0xFE | "HELLO, WORLD" |
 +------+--------+------+----------------+      +------+----------------+
```
- `粘包`和`拆包`：
  - 业务中使用 TCP 协议获得一个有序的稳定的流作为消息的载体。
  - 但在实践中应用层不断向传输层轮循缓冲区的数据的时候会得到片段化的数据，
  - 虽然有序但是这些缓冲中的数据可能归属是多个不同的消息（造词：粘包）。
  - 将这些数据以特定形式（例如包长度(length)+负载(payload)）拆分或重组成逻辑消息（造词：拆包）
  
- 实现一个从 socket connection 中解码出 goim 协议的解码器。
  - goim over `UDP`   
    - server
    ```bash
    go run goim_udp.go server
    ```
    - client
    ```
    go run goim_udp.go client -msg auth
    ```