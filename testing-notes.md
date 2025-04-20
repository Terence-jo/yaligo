## Testing the Parser (readFromTokens)
It works really well, but there seems to be trailing junk on the resulting linked list. It 
matches the reference list to the end, but then the result of the parser looks like it has 
and additional '+' which is perplexing. Is there anywhere it might double-parse a token, especially
if it enters and then exits a list? The plus is the first element of the inner list though, so that
is very strange.
    - That's all resolved now. The parsing works, I had an issue keeping track of position when recursing.

## What now?
- The parsing seems to work exactly as intended, so now I should look to evaluation.
- I can construct an internal representation, now I need to be able to compute the results of expressions
captured in that representation.
