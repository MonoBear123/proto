import os
import yaml
from typing import Any


class Config:
    def __init__(self, config_file_path: str):
        self.config_file_path = config_file_path
        self._config = None

    def load_config(self) -> None:
        if not os.path.exists(self.config_file_path):
            raise FileNotFoundError(f"Configuration file {self.config_file_path} not found.")

        with open(self.config_file_path, 'r', encoding='utf-8') as file:
            try:
                self._config = yaml.safe_load(file)
            except yaml.YAMLError as e:
                raise RuntimeError(f"Error parsing YAML file: {e}")

    def get(self, key: str, default: Any = None) -> Any:
        if self._config is None:
            self.load_config()

        keys = key.split('.')
        value = self._config
        for k in keys:
            if isinstance(value, dict):
                value = value.get(k, default)
            else:
                return default
        return value


