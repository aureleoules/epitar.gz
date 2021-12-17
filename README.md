# epitar.gz
Highly customizable archive and index framework for EPITA.

## How does it work

### Archive modules
An archive module scrapes, downloads, or archives websites and services. These modules are highly customizable as they run in Docker containers.

### Index
Archived files may be scanned to build a search index.
PDF files words are extracted using regular methods or using an OCR for scanned documents.  
Words are then processed by a [sonic](https://github.com/valeriansaliou/sonic) instance in order to build a fast search index.

### UI & API
A UI is exposed along with an API to quickly search for files.

## Contributing

### Add an archive module

An archive module is highly customizable as it can be written in programming language as long as a valid `Dockerfile` is provided.  
Your archive module **must** have a `Dockerfile`, a `module.json` and a `README`.

#### Dockerfile
Your `Dockerfile` **can** use any base image but try to keep the image size small.

The output directory for archived files **must** be `/output`.

#### module.json
Your `module.json` **must** provide informations about the website or service that is being archived.  
Here is an example:  
```json
{
    "name": "Past-Exams",
    "slug": "past-exams",
    "url": "https://github.com/Epidocs/Past-Exams",
    "description": "Past subjects and other files, for the benefit of EPITA students. ",
    "logo": "https://github.com/fluidicon.png", // optional
    "authors": [
        {
            "name": "Aurele Oules",
            "email": "aurele@oules.com"
        }
    ]
}
```

#### README.md
You **must** provide a simple `README.md` that explains how to use this module.  
An archive module **may** take environment variables as options so you may explain them here.

#### Other files
You **may** add any other files in the module directory but try to keep it organized and only commit necessary files.

You **must** edit the `config.sample.yml` file to provide an example on how to use your archive module.

## License
[MIT](https://github.com/aureleoules/bitcandle/blob/master/LICENSE) - [Aurèle Oulès](https://www.aureleoules.com)
