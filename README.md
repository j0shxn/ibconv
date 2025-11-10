# ibconv

ITU e-bulletin easy image converter.

This is a simple bulk image converter. The program converts all images (jpg,
png, gif) from an input folder, resizes them, and saves them to an output
folder in the specified format.

    Default behavior (no arguments):
      Converts images from ./source to ./sink at 280x180 resolution
      in 'jpg' format.

    Options:
      -i <folder>     Path to the input source folder.
                      (default: ./source)

      -o <folder>     Path to the output sink folder.
                      (default: ./sink)

      -r <W,H>        Target resolution (Width,Height) to resize images to.
                      (e.g., "800,600")
                      (default: 280,180)

      -f <format>     Target output format. Can be 'jpg' or 'png'.
                      (default: jpg)

      -v              Enable verbose output, showing processing details.

      -h              Print this help and usage message.
