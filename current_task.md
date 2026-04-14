先定 Lexer 状态：input / position / readPosition / ch / line / column，采用“单字符前瞻”模型。
实现 3 个底层移动函数：readChar()（前进）、peekChar()（只看不动）、newToken()（统一生成 token 并带位置信息）。
实现 4 个读取函数：skipWhitespace()、skipComment()（先支持 //）、readIdentifier()、readNumber()、readString()。
最后实现 NextToken()：先跳空白和注释，再按 switch 识别单/双字符运算符（= vs ==、! vs !=），标识符走关键字表，数字和字符串走专用读取器。
关键细节：token 的 Line/Column 记录“起始位置”，不是读完后的位置。

