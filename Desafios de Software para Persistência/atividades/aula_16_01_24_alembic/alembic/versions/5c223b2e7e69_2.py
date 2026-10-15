"""2

Revision ID: 5c223b2e7e69
Revises: 7adba4906f43
Create Date: 2025-01-16 09:56:59.398478

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = '5c223b2e7e69'
down_revision: Union[str, None] = '7adba4906f43'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    pass


def downgrade() -> None:
    pass
