# JSSP

#### Jerry's Super Simple Protocol


Toy protocol to Serialize numbers or strings to a byte stream,
using a length-value representation.


jprot: jerry protocol

**NUMBER**		{+-}?{0..9}+		Max. length   9 digits + 1 for sign

**TEXTSTRING**	{Unicode char}*		Max. length 256 bytes

**MESSAGE**		{<NUMBER> | <TEXTSTRING>}+

Encoding uses LENGTH-CONTENT encoding:
LENGTH	1 byte, value is the number of content bytes
CONTENT	0 or more content bytes

Example:
```
"hello" -> []byte{ 5 , h , e , l , l , o}
42      -> []byte{ 2, '4', '2'}
```
