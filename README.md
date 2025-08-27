# CBZ to PDF

Simple and quick way to convert CBZ files to PDF files without losing quality.

## Demo

Simply drag and drop the executable into the <ins>same folder with your cbz files</ins>. Double click on the executable and PDF files will be created in the same folder (one `.pdf` for each `.cbz` file)



## Downloading

### Download executable (Recommended)

Download a [binary for your platform](https://github.com/Serious-Fin/type-training/releases/tag/v1.0.0) from release page.

On linux/mac you might need to add execution permission to file:

```zsh
chmod +x <binary name>
```

### Build from source (Advanced)

Open a console window and clone repository locally:

```zsh
git clone https://github.com/Serious-Fin/type-training.git
```

Move into the cloned repository

```zsh
cd type-training
```

Install all go dependencies (you might need to [install go if you haven't already](https://go.dev/doc/install) before this step)

```zsh
go install
```

Build the project

```zsh
go build
```

Now you should have an executable `type-training`. Run it with:

```zsh
./type-training
```

