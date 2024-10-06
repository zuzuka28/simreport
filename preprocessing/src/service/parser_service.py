import hashlib
import io
import os
import re
import xml.etree.ElementTree as ET
import zipfile
from dataclasses import dataclass

from src.domain.model.document import Document
from src.domain.model.image import Image
from src.domain.service.parser_service import DocumentParserService


@dataclass
class InDocxImage:
    fname: str
    data: bytes


# code snippet from <https://github.com/ankushshah89/python-docx2txt/blob/master/docx2txt/docx2txt.py>
class ParsedDocx:
    _nsmap = {"w": "http://schemas.openxmlformats.org/wordprocessingml/2006/main"}

    def __init__(self, docx: bytes) -> None:
        self._source = docx

        self._text = ""
        self._images: list[InDocxImage] = []

        self._process()

    def get_text(self) -> str:
        return self._text

    def get_images(self) -> list[InDocxImage]:
        return self._images

    def _qn(self, tag):
        """
        Stands for 'qualified name', a utility function to turn a namespace
        prefixed tag name into a Clark-notation qualified tag name for lxml. For
        example, ``qn('p:cSld')`` returns ``'{http://schemas.../main}cSld'``.
        Source: https://github.com/python-openxml/python-docx/
        """
        prefix, tagroot = tag.split(":")
        uri = self._nsmap[prefix]
        return "{{{}}}{}".format(uri, tagroot)

    def _xml2text(self, xml):
        """
        A string representing the textual content of this run, with content
        child elements like ``<w:tab/>`` translated to their Python
        equivalent.
        Adapted from: https://github.com/python-openxml/python-docx/
        """
        text = ""
        root = ET.fromstring(xml)
        for child in root.iter():
            if child.tag == self._qn("w:t"):
                t_text = child.text
                text += t_text if t_text is not None else ""
            elif child.tag == self._qn("w:tab"):
                text += "\t"
            elif child.tag in (self._qn("w:br"), self._qn("w:cr")):
                text += "\n"
            elif child.tag == self._qn("w:p"):
                text += "\n\n"
        return text

    # TODO: extract other types of embedded files
    # TODO: extract tables
    # TODO: header and footer processing
    def _process(
        self,
        add_header=False,
        add_footer=False,
    ):
        text = ""

        # unzip the docx in memory
        with zipfile.ZipFile(io.BytesIO(self._source)) as zipf:
            filelist = zipf.namelist()

            if add_header:
                # get header text
                # there can be 3 header files in the zip
                header_xmls = "word/header[0-9]*.xml"
                for fname in filelist:
                    if re.match(header_xmls, fname):
                        text += self._xml2text(zipf.read(fname))

            # get main text
            doc_xml = "word/document.xml"
            text += self._xml2text(zipf.read(doc_xml))

            if add_footer:
                # get footer text
                # there can be 3 footer files in the zip
                footer_xmls = "word/footer[0-9]*.xml"
                for fname in filelist:
                    if re.match(footer_xmls, fname):
                        print(zipf.read(fname))
                        text += self._xml2text(zipf.read(fname))

            self._text = text

            for fname in filelist:
                _, extension = os.path.splitext(fname)
                if extension in [".jpg", ".jpeg", ".png", ".bmp"]:
                    self._images.append(
                        InDocxImage(
                            **{
                                "fname": os.path.basename(fname),
                                "data": zipf.read(fname),
                            }
                        )
                    )


class ParserService(DocumentParserService):
    def __init__(self):
        pass

    async def parse_document(self, source: bytes) -> Document:
        sha256 = _sha256_bytes(source)

        parsed = ParsedDocx(source)

        return Document(
            **{
                "id": sha256,
                "source_bytes": source,
                "sha256": sha256,
                "text_content": parsed.get_text(),
                "images": self._parsed_docx_images(parsed),
            }
        )

    @staticmethod
    def _parsed_docx_images(source: ParsedDocx) -> list[Image]:
        raw_images = source.get_images()

        return [
            Image(
                **{
                    "fname": img.fname,
                    "source_bytes": img.data,
                    "sha256": _sha256_bytes(img.data),
                }
            )
            for img in raw_images
        ]


def _sha256_bytes(source: bytes) -> str:
    sha256_hash = hashlib.sha256()
    sha256_hash.update(source)
    return sha256_hash.hexdigest()
