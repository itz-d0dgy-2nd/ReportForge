---
FindingID: Web-01
FindingName: web_authn_broken_authentication
FindingTitle: Insecure Authentication
FindingStatus: Unresolved
FindingImpact: Severe
FindingLikelihood: Possible
FindingAuthor: JM
FindingReviewers: 
---

### Retesting 2

<retest_fixed>

- Client Variable: !Client
- Target Variable: !TargetAsset0 || !TargetAsset1
  -  First un-ordered Subitem
  - Second un-ordered Subitem 
- Third un-ordered item

</retest_fixed>

### Retesting 1

<retest_not_fixed>

- First un-ordered item
- Second un-ordered item
  - First un-ordered Subitem
  - Second un-ordered Subitem 
- Third un-ordered item

</retest_not_fixed>

### Details

This is a paragraph with **bold text**, *italic text*, <mark>highlighted</mark> text, `inline code`, and <qa>QA TEXT</qa>. I am now referencing another finding [here](#Critical_mobile_authn_broken_mfa). This sentance was added to test report gen in action

- First un-ordered item
- Second un-ordered item
  -  First un-ordered Subitem
  - Second un-ordered Subitem 
- Third un-ordered item

1. First ordered item
1. Second ordered item
1. Third ordered item

[Google](https://www.google.com)

![This is the alt text for the troll img](report/Screenshots/troll.png)

> This is a blockquote.

text 

```php
// This is a PHP code block
<?php echo "Hello World"; ?>
```

text 

```csharp
// This is a csharp code block
Console.WriteLine("Hello World");
```
text 

```python
# This is a python code block
print("Hello World")
```
text 
```go
// This is a golang code block
fmt.Println("Hello World")
```

| Name | Age |
| ---- | --- |
| Bob  | 27  |
| Jo   | 23  |

### Reproduction Steps

### Recommendations
