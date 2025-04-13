import random
from typing import List, Set, Dict
import string

import nltk
from nltk.tokenize import sent_tokenize, word_tokenize
from nltk.corpus import words as nltk_words


nltk.download("punkt")
nltk.download("punkt_tab")
nltk.download("words")


class SentenceAugmentor:
    def __init__(self, language: str = "en") -> None:
        """
        Инициализирует объект SentenceAugmentor.

        :param language: Язык текста (по умолчанию "en").
        """
        self.language = language
        self.prefix = "In addition, "
        self.suffix = " Furthermore, "
        self.inserted_sentence = "This is a randomly inserted sentence."
        self.substituted_sentence = "This is a substituted sentence."

    def split_into_sentences(self, text: str) -> List[str]:
        """
        Разбивает текст на предложения.

        :param text: Входной текст.
        :return: Список предложений.
        """
        return sent_tokenize(text, language=self.language)

    def random_deletion(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Удаляет случайные предложения из текста.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        sentences = self.split_into_sentences(text)
        if not sentences or len(sentences) <= 1:
            return text

        num_to_delete = max(int(len(sentences) * augmentation_rate), 1)
        for _ in range(num_to_delete):
            if len(sentences) == 0:
                break
            del_idx = random.randint(0, len(sentences) - 1)
            del sentences[del_idx]
        return " ".join(sentences)

    def random_truncation(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Усекает случайные предложения в тексте.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        sentences = self.split_into_sentences(text)
        if not sentences:
            return text
        num_to_truncate = int(len(sentences) * augmentation_rate)
        for _ in range(num_to_truncate):
            trunc_idx = random.randint(0, len(sentences) - 1)
            sentences[trunc_idx] = sentences[trunc_idx][
                : random.randint(1, len(sentences[trunc_idx]) - 1)
            ]
        return " ".join(sentences)

    def random_prefix(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Добавляет случайный префикс к случайным предложениям в тексте.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        sentences = self.split_into_sentences(text)
        if not sentences:
            return text
        num_to_prefix = int(len(sentences) * augmentation_rate)
        for _ in range(num_to_prefix):
            idx = random.randint(0, len(sentences) - 1)
            sentences[idx] = self.prefix + sentences[idx]
        return " ".join(sentences)

    def random_suffix(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Добавляет случайный суффикс к случайным предложениям в тексте.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        sentences = self.split_into_sentences(text)
        if not sentences:
            return text
        num_to_suffix = int(len(sentences) * augmentation_rate)
        for _ in range(num_to_suffix):
            idx = random.randint(0, len(sentences) - 1)
            sentences[idx] = sentences[idx] + self.suffix
        return " ".join(sentences)

    def random_insertion(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Вставляет случайные предложения в случайные позиции в тексте.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        sentences = self.split_into_sentences(text)
        num_to_insert = int(len(sentences) * augmentation_rate)
        for _ in range(num_to_insert):
            insert_idx = random.randint(0, len(sentences))
            sentences.insert(insert_idx, self.inserted_sentence)
        return " ".join(sentences)

    def repeat_sentence(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Повторяет случайные предложения в тексте.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        sentences = self.split_into_sentences(text)
        if not sentences:
            return text
        num_to_repeat = max(int(len(sentences) * augmentation_rate), 1)
        for _ in range(num_to_repeat):
            repeat_idx = random.randint(0, len(sentences) - 1)
            sentences.append(sentences[repeat_idx])
        return " ".join(sentences)

    def lowercase_uppercase_sentence(
        self, text: str, augmentation_rate: float = 0.1
    ) -> str:
        """
        Преобразует случайные предложения в нижний или верхний регистр.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        sentences = self.split_into_sentences(text)
        if not sentences:
            return text
        num_to_change = int(len(sentences) * augmentation_rate)
        for _ in range(num_to_change):
            idx = random.randint(0, len(sentences) - 1)
            if sentences[idx].islower():
                sentences[idx] = sentences[idx].upper()
            else:
                sentences[idx] = sentences[idx].lower()
        return " ".join(sentences)

    def random_substitution(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Заменяет случайные предложения на заданное.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        sentences = self.split_into_sentences(text)
        if not sentences:
            return text
        num_to_substitute = int(len(sentences) * augmentation_rate)
        for _ in range(num_to_substitute):
            sub_idx = random.randint(0, len(sentences) - 1)
            sentences[sub_idx] = self.substituted_sentence
        return " ".join(sentences)

    def neighboring_swap(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Меняет местами соседние предложения в тексте.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        sentences = self.split_into_sentences(text)
        if len(sentences) < 2:
            return text
        num_to_swap = max(int(len(sentences) * augmentation_rate), 1)
        for _ in range(num_to_swap):
            idx = random.randint(0, len(sentences) - 2)
            sentences[idx], sentences[idx + 1] = sentences[idx + 1], sentences[idx]
        return " ".join(sentences)


class WordAugmentor:
    def __init__(self, language: str = "en") -> None:
        """
        Инициализация класса для модификации слов.

        :param language: Язык текста (по умолчанию 'en' для английского).
        """
        self.language = language

        self.english_words: Set[str] = set(nltk_words.words())

        self.words = self.english_words

        self.punctuation: str = string.punctuation

    def split_into_words(self, text: str) -> List[str]:
        """
        Разбивает текст на слова.

        :param text: Входной текст.
        :return: Список слов.
        """
        return word_tokenize(text)

    def random_deletion(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Удаляет случайные слова из текста.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        words: List[str] = self.split_into_words(text)
        if len(words) == 0:
            return text
        num_to_delete: int = int(len(words) * augmentation_rate)
        for _ in range(num_to_delete):
            if len(words) == 0:
                break
            del_idx: int = random.randint(0, len(words) - 1)
            del words[del_idx]
        return " ".join(words)

    def random_insertion(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Вставляет случайные слова в случайные позиции текста.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        words: List[str] = self.split_into_words(text)
        if not words:
            return text

        num_to_insert: int = int(len(words) * augmentation_rate)
        if num_to_insert == 0:
            return text

        new_words: List[str] = random.choices(list(self.words), k=num_to_insert)

        insert_indices: List[int] = sorted(
            random.sample(range(len(words) + 1), num_to_insert)
        )

        for idx, word in zip(insert_indices, new_words):
            words.insert(idx, word)

        return " ".join(words)

    def random_substitution(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Заменяет случайные слова на случайные другие слова.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        words: List[str] = self.split_into_words(text)
        if not words:
            return text

        num_to_substitute: int = int(len(words) * augmentation_rate)
        if num_to_substitute == 0:
            return text

        sub_indices = random.sample(range(len(words)), num_to_substitute)

        new_words = random.choices(list(self.words), k=num_to_substitute)

        for idx, new_word in zip(sub_indices, new_words):
            words[idx] = new_word

        return " ".join(words)

    def repeat_word(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Повторяет случайные слова в тексте.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        words: List[str] = self.split_into_words(text)
        if len(words) == 0:
            return text
        num_to_repeat: int = int(len(words) * augmentation_rate)
        for _ in range(num_to_repeat):
            repeat_idx: int = random.randint(0, len(words) - 1)
            words[repeat_idx] = words[repeat_idx] * 2  # Повторяем слово дважды
        return " ".join(words)

    def n_gram_based_substitution(
        self, text: str, n: int = 3, augmentation_rate: float = 0.1
    ) -> str:
        """
        Заменяет слова на основе n-граммной частотности.

        :param text: Входной текст.
        :param n: Размер n-граммы.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        words: List[str] = self.split_into_words(text)
        if len(words) < n:
            return text
        num_to_change: int = int(len(words) * augmentation_rate)
        for _ in range(num_to_change):
            idx: int = random.randint(n - 1, len(words) - 1)
            new_word: str = random.choice(list(self.words))
            words[idx] = new_word
        return " ".join(words)

    def neighboring_swap(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Меняет местами соседние слова.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        words: List[str] = self.split_into_words(text)
        if len(words) < 2:
            return text
        num_to_swap: int = int(len(words) * augmentation_rate)
        for _ in range(num_to_swap):
            idx: int = random.randint(0, len(words) - 2)
            words[idx], words[idx + 1] = words[idx + 1], words[idx]
        return " ".join(words)


class CharacterAugmentor:
    def __init__(self, language: str = "en") -> None:
        """
        Инициализация класса для модификации символов.

        :param language: Язык текста (по умолчанию 'en' для английского).
        """
        self.language = language
        self.punctuation: str = string.punctuation
        self.ascii_letters: str = string.ascii_letters
        self.ascii_digits: str = string.digits
        self.language_alphabets: Dict[str, str] = {
            "en": string.ascii_lowercase + string.ascii_uppercase,
            "ru": "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ",
        }
        self.qwerty_map: Dict[str, str] = {
            "q": "w",
            "w": "e",
            "e": "r",
            "r": "t",
            "t": "y",
            "y": "u",
            "u": "i",
            "i": "o",
            "o": "p",
            "a": "s",
            "s": "d",
            "d": "f",
            "f": "g",
            "g": "h",
            "h": "j",
            "j": "k",
            "k": "l",
            "z": "x",
            "x": "c",
            "c": "v",
            "v": "b",
            "b": "n",
            "n": "m",
            "Q": "W",
            "W": "E",
            "E": "R",
            "R": "T",
            "T": "Y",
            "Y": "U",
            "U": "I",
            "I": "O",
            "O": "P",
            "A": "S",
            "S": "D",
            "D": "F",
            "F": "G",
            "G": "H",
            "H": "J",
            "J": "K",
            "K": "L",
            "Z": "X",
            "X": "C",
            "C": "V",
            "V": "B",
            "B": "N",
            "N": "M",
        }
        self.qwerty_ru_map = {
            "й": "ц",
            "ц": "й",
            "у": "к",
            "к": "у",
            "е": "н",
            "н": "е",
            "г": "ш",
            "ш": "г",
            "щ": "з",
            "з": "щ",
            "х": "ъ",
            "ъ": "х",
            "ф": "ы",
            "ы": "ф",
            "в": "в",
            "а": "а",
            "п": "р",
            "р": "п",
            "о": "о",
            "л": "д",
            "д": "л",
            "ж": "э",
            "э": "ж",
            "я": "я",
            "ч": "с",
            "с": "ч",
            "м": "ь",
            "ь": "м",
            "и": "и",
            "т": "т",
            "ь": "м",
            "м": "ь",
            "б": "ю",
            "ю": "б",
            "ю": "б",
            "б": "ю",
            "ё": "ё",
            "Й": "Ц",
            "Ц": "Й",
            "У": "К",
            "К": "У",
            "Е": "Н",
            "Н": "Е",
            "Г": "Ш",
            "Ш": "Г",
            "Щ": "З",
            "З": "Щ",
            "Х": "Ъ",
            "Ъ": "Х",
            "Ф": "Ы",
            "Ы": "Ф",
            "В": "В",
            "А": "А",
            "П": "Р",
            "Р": "П",
            "О": "О",
            "Л": "Д",
            "Д": "Л",
            "Ж": "Э",
            "Э": "Ж",
            "Я": "Я",
            "Ч": "С",
            "С": "Ч",
            "М": "Ь",
            "Ь": "М",
            "И": "И",
            "Т": "Т",
            "Ь": "М",
            "М": "Ь",
            "Б": "Ю",
            "Ю": "Б",
            "Ю": "Б",
            "Б": "Ю",
            "Ё": "Ё",
        }

        if language == "russian":
            self.qwerty_map = self.qwerty_ru_map

        self.homoglyphs: Dict[str, str] = {
            "a": "ɑ",
            "A": "А",
            "e": "е",
            "E": "Е",
            "o": "о",
            "O": "О",
            "c": "с",
            "C": "С",
        }

    def random_deletion(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Удаляет случайные символы из текста.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        if not text:
            return text
        text_list: List[str] = list(text)
        num_to_delete: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_delete):
            if len(text_list) == 0:
                break
            del_idx: int = random.randint(0, len(text_list) - 1)
            del text_list[del_idx]
        return "".join(text_list)

    def case_substitution(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Заменяет регистр символов в тексте.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_change: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_change):
            idx: int = random.randint(0, len(text_list) - 1)
            if text_list[idx].islower():
                text_list[idx] = text_list[idx].upper()
            else:
                text_list[idx] = text_list[idx].lower()
        return "".join(text_list)

    def n_gram_based_substitution(
        self, text: str, n: int = 3, augmentation_rate: float = 0.1
    ) -> str:
        """
        Заменяет символы на основе n-граммной частотности.

        :param text: Входной текст.
        :param n: Размер n-граммы.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_change: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_change):
            idx: int = random.randint(0, len(text_list) - n)
            for j in range(n):
                if random.random() < 0.5:
                    text_list[idx + j] = random.choice(
                        self.ascii_letters + self.ascii_digits
                    )
        return "".join(text_list)

    def qwerty_typo_substitution(
        self, text: str, augmentation_rate: float = 0.1
    ) -> str:
        """
        Заменяет символы на основе карты QWERTY для имитации опечаток.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_change: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_change):
            idx: int = random.randint(0, len(text_list) - 1)
            if text_list[idx] in self.qwerty_map:
                text_list[idx] = self.qwerty_map[text_list[idx]]
        return "".join(text_list)

    def homoglyph_substitution(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Заменяет символы на их гомоглифы.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_change: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_change):
            idx: int = random.randint(0, len(text_list) - 1)
            if text_list[idx] in self.homoglyphs:
                text_list[idx] = self.homoglyphs[text_list[idx]]
        return "".join(text_list)

    def random_ascii_substitution(
        self, text: str, augmentation_rate: float = 0.1
    ) -> str:
        """
        Заменяет случайные символы на случайные ASCII-символы.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_change: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_change):
            idx: int = random.randint(0, len(text_list) - 1)
            if text_list[idx] in self.punctuation:
                text_list[idx] = random.choice(self.ascii_letters + self.ascii_digits)
        return "".join(text_list)

    def random_character_from_language_alphabet_substitution(
        self, text: str, augmentation_rate: float = 0.1
    ) -> str:
        """
        Заменяет случайные символы на символы из алфавита указанного языка.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        language_alphabet: str = self.language_alphabets.get(
            self.language, string.ascii_letters
        )
        num_to_change: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_change):
            idx: int = random.randint(0, len(text_list) - 1)
            text_list[idx] = random.choice(language_alphabet)
        return "".join(text_list)

    def random_punctuation_substitution(
        self, text: str, augmentation_rate: float = 0.1
    ) -> str:
        """
        Заменяет случайные символы на знаки препинания.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_change: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_change):
            idx: int = random.randint(0, len(text_list) - 1)
            if text_list[idx] in self.punctuation:
                text_list[idx] = random.choice(self.punctuation)
        return "".join(text_list)

    def random_unicode_character_substitution(
        self, text: str, augmentation_rate: float = 0.1
    ) -> str:
        """
        Заменяет случайные символы на случайные символы Unicode.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_change: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_change):
            idx: int = random.randint(0, len(text_list) - 1)
            text_list[idx] = chr(random.randint(0, 0x10FFFF))
        return "".join(text_list)

    def character_repetition(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Повторяет случайные символы в тексте.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        if not text_list:
            return text
        num_to_repeat: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_repeat):
            idx: int = random.randint(0, len(text_list) - 1)
            text_list[idx] = text_list[idx] * (
                random.randint(2, 5)
            )  # Repeat 2 to 5 times
        return "".join(text_list)

    def n_grams_based_insertion(
        self, text: str, n: int = 3, augmentation_rate: float = 0.1
    ) -> str:
        """
        Вставляет случайные n-граммы в текст.

        :param text: Входной текст.
        :param n: Размер n-граммы.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_insert: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_insert):
            pos: int = random.randint(0, len(text_list))
            new_n_gram: str = "".join(
                random.choices(string.ascii_letters + string.digits, k=n)
            )
            text_list.insert(pos, new_n_gram)
        return "".join(text_list)

    def random_character_from_language_alphabet_insertion(
        self, text: str, augmentation_rate: float = 0.1
    ) -> str:
        """
        Вставляет случайные символы из алфавита указанного языка в текст.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        language_alphabet: str = self.language_alphabets.get(
            self.language, string.ascii_letters
        )
        num_to_insert: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_insert):
            pos: int = random.randint(0, len(text_list))
            new_char: str = random.choice(language_alphabet)
            text_list.insert(pos, new_char)
        return "".join(text_list)

    def random_punctuation_insertion(
        self, text: str, augmentation_rate: float = 0.1
    ) -> str:
        """
        Вставляет случайные знаки препинания в текст.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_insert: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_insert):
            pos: int = random.randint(0, len(text_list))
            new_punct: str = random.choice(self.punctuation)
            text_list.insert(pos, new_punct)
        return "".join(text_list)

    def random_unicode_character_insertion(
        self, text: str, augmentation_rate: float = 0.1
    ) -> str:
        """
        Вставляет случайные символы Unicode в текст.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        text_list: List[str] = list(text)
        num_to_insert: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_insert):
            pos: int = random.randint(0, len(text_list))
            new_unicode_char: str = chr(random.randint(0, 0x10FFFF))
            text_list.insert(pos, new_unicode_char)
        return "".join(text_list)

    def neighboring_swap(self, text: str, augmentation_rate: float = 0.1) -> str:
        """
        Меняет местами два соседних символа.

        :param text: Входной текст.
        :param augmentation_rate: Процент аугментации (от 0 до 1).
        :return: Измененный текст.
        """
        if len(text) < 2:
            return text
        text_list: List[str] = list(text)
        num_to_swap: int = int(len(text_list) * augmentation_rate)
        for _ in range(num_to_swap):
            idx: int = random.randint(0, len(text_list) - 2)
            text_list[idx], text_list[idx + 1] = text_list[idx + 1], text_list[idx]
        return "".join(text_list)
