# mango

*Generate manual pages from the source code of your Go commands*

*mango* is a small command line utility that allows you to create manual
pages from the source code of your Go commands.

## Overview

TODO: ...

## Formatting

mango supports comment formatting in a Markdown-like syntax.

### Headings

```
First Heading:

Riverrun, past Eve and Adam's, from swerve of shore to bend
of bay, brings us by a commodius vicus of recirculation
back to Howth Castle and Environs.

Second Heading
==============
Riverrun, past Eve and Adam's, from swerve of shore to bend
of bay, brings us by a commodius vicus of recirculation
back to Howth Castle and Environs.

Third Heading
-------------
Riverrun, past Eve and Adam's, from swerve of shore to bend
of bay, brings us by a commodius vicus of recirculation
back to Howth Castle and Environs.
```

### Paragraphs

Paragraphs are separated by a blank line.

```
Riverrun, past Eve and Adam's, from swerve of shore to bend
of bay, brings us by a commodius vicus of recirculation
back to Howth Castle and Environs.

Sir Tristram, violer d'amores, fr'over the short sea,
had passencore rearrived from North Armorica on this side the scraggy
isthmus of Europe Minor to wielderfight his penisolate war.
```

### Emphasis

Asterisks (\*) and underscores (\_) are indicators of emphasis.

```
*Wassaily Booslaeugh* of _Riesengeborg_
```

### Code

Code blocks begin with a closing angle bracket (>).

```
> echo "Kick nuck, Knockcastle"
```

### Lists

Asterisks, numbers or single words followed by a closing parenthesis
at the start of a line create list items.

```
*) Item
*) Item
*) Item

1) Item
2) Item
3) Item

a) Item
b) Item
c) Item
```

### License

mango is released under MIT license.
You can find a copy of the MIT License in the [LICENSE](./LICENSE) file.

