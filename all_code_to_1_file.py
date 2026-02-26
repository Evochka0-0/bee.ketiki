import os

start_dir = os.getcwd()
output_file = "output.txt"

with open(output_file, "w", encoding="utf-8") as out_f:
    for root, dirs, files in os.walk(start_dir):
        for file_name in files:
            # Пропускаем файл output.txt, чтобы избежать дублирования
            if file_name == output_file:
                continue
                
            full_path = os.path.join(root, file_name)
            rel_path = os.path.relpath(full_path, start_dir)
            out_f.write(f"полный путь от исходной папки: {rel_path}\n")
            try:
                with open(full_path, "r", encoding="utf-8") as f:
                    content = f.read()
                out_f.write("```\n")
                out_f.write(content)
                if not content.endswith("\n"):
                    out_f.write("\n")
                out_f.write("```\n\n")
            except Exception as e:
                folder = os.path.dirname(rel_path)
                folder = folder if folder else "исходная папка"
                out_f.write(
                    f"В папке {folder} лежит файл {file_name}, который не удалось прочитать.\n\n"
                )