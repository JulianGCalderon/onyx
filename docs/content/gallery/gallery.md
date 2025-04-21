# The Rock Gallery

This page showcases examples of the Markdown features supported by Onyx.

## Text Formatting

You can apply **bold** using either `__` or `**`.

You can apply _italic_ using either `_` or `*`.

Use double tilde `~~` for ~~strikethrough text~~.

Use backticks `` ` `` for `inline code`.

Here’s an [external link](https://commonmark.org/help/).

And here’s an autolink https://commonmark.org/help/.

## Blockquotes

Blockquotes are used for emphasized content.

> This is a blockquote.

They are also used for actual quotes.

> "Mientras uno está vivo, uno debe amar lo más que pueda"
>
> Jacobo Morales

## Lists

Use hyphens (`-`) or asterisks (`*`) to create unordered lists.

- This is an unordered list.
- It can contain many elements.
- As many as you want.

Use numbers followed by a period (`1.`) to create ordered lists.

1. Lists can also be ordered.
2. They support multiple elements.
3. As many as needed.

## Code Blocks

Use triple backticks `` ``` `` to display blocks of code in a readable format.

```c
int main() {
  return 0;
}
```

Inside of code blocks, there is no line wrapping:

```c
int main(int argc, char** argv) {
        if (argc < 9) {
                if (argc < 8) {
                        if (argc < 7) {
                                if (argc < 6) {
                                        if (argc < 5) {
                                                if (argc < 4) {
                                                        return 1;
                                                }
                                        }
                                }
                        }
                }
        }
        return 0;
}
```

## Footnotes

Footnotes allow you to include additional information without cluttering the main text.

Here's a simple footnote[^1], and here's a longer one with more detail[^bignote].

## Tables

Tables are great for displaying structured data in rows and columns.

| Lorem      | Ipsum       | Dolor      | Sit    |
|------------|-------------|------------|--------|
| Amet       | Consectetur | Adipiscing | Elit   |
| Sed        | Do          | Eiusmod    | Tempor |
| Incididunt | Ut          | Labore     | Et     |

## Definition Lists

Definition lists let you associate terms with one or more definitions.

First Term
: This is the definition of the first term.

Second Term
: This is one definition of the second term.
: This is another definition of the second term.

## Nested Blocks

Many block elements have for support nested blocks. For example: blockquotes.

Lists can contain other lists inside.

1. Fruit: The sweet an juicy product of a tree or other plant that contains seed.

   Some example of fruits are:
   - Apples
   - Oranges
   - Bananas

2. Vegetables: A plant or part of a plant used as food.
   - Cabbage
   - Potato
   - Carrot

Definition lists can contain lists inside.

Fruit
: The sweet an juicy product of a tree or other plant that contains seed.
  - Apples
  - Oranges
  - Bananas

Vegetables
: A plant or part of a plant used as food.
  - Cabbage
  - Potato
  - Carrot

Blockquotes can contain lists inside.

> A Fruit is the sweet an juicy product of a tree or other plant that contains seed.
>
> - Apples
> - Oranges
> - Bananas
>
> ~ Oxford Languages

Lists can contain code blocks inside.

- Rust
  ```rust
  fn main() {
    println!("Hello World");
  }
  ```
- Python
  ```python
  def main():
    print("Hello, world!")
  ```

## Wikilinks

Wikilinks are internal links between notes, commonly used for easy navigation.

- [[gallery/inner/absolute]]: An absolute path to a note.
- [[inner/relative]]: A relative path to a note within the current folder.
- [[unique]]: A uniquely named note that can be linked without a path.
- [[gallery/inner/absolute|absolute]]: An absolute path to a note, with a different title.
- [[#Nested Blocks|nested blocks]]: A path to a heading.

Want to go back to the [[index]]?

[^1]: This is the first footnote.

[^bignote]: This is a longer footnote that can include more context.

    It can span multiple lines, or even contain blocks.
