import os
import logging
from pathlib import Path
from huggingface_hub import snapshot_download

snapshot_download(
    repo_id='Code0871/measure_buyer',
    local_dir='./ai_model',
    local_dir_use_symlinks=False,
    resume_download=True
)