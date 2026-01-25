import sys
import os

try:
    import pdfplumber
except ImportError:
    print("pdfplumber not installed. Please run 'pip install pdfplumber'")
    sys.exit(1)


def read_pdf(file_path):
    """
    读取 PDF 文件并提取每页的文本内容
    """
    if not os.path.exists(file_path):
        print(f"Error: File not found at {file_path}")
        return

    print(f"Reading PDF from: {file_path}\n" + "-" * 30)

    try:
        full_text = ""
        with pdfplumber.open(file_path) as pdf:
            for i, page in enumerate(pdf.pages):
                text = page.extract_text()
                if text:
                    page_content = f"\n--- Page {i+1} ---\n{text}\n"
                    print(page_content)
                    full_text += page_content

        return full_text
    except Exception as e:
        print(f"An error occurred: {e}")


if __name__ == "__main__":
    # 默认路径
    target_file = ""

    # 允许命令行覆盖
    if len(sys.argv) > 1:
        target_file = sys.argv[1]

    read_pdf(target_file)
