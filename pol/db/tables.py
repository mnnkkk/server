"""
this file is mostly generated by sqlacodegen from database
"""
from typing import Type

from sqlalchemy import Enum, Index, Table, Column, String, text
from sqlalchemy.dialects.mysql import (
    ENUM,
    YEAR,
    INTEGER,
    TINYINT,
    VARCHAR,
    SMALLINT,
    MEDIUMINT,
    MEDIUMTEXT,
)
from sqlalchemy.ext.declarative import DeclarativeMeta, declarative_base

Base: Type[DeclarativeMeta] = declarative_base()
metadata = Base.metadata


class Character(Base):
    __tablename__ = "chii_characters"

    id = Column("crt_id", MEDIUMINT(8), primary_key=True)
    name = Column("crt_name", String(255, "utf8_unicode_ci"), nullable=False)
    crt_role = Column(TINYINT(4), nullable=False, index=True, comment="角色，机体，组织。。")
    crt_infobox = Column(MEDIUMTEXT, nullable=False)
    crt_summary = Column(MEDIUMTEXT, nullable=False)
    crt_img = Column(String(255, "utf8_unicode_ci"), nullable=False)
    crt_comment = Column(MEDIUMINT(9), nullable=False, server_default=text("'0'"))
    crt_collects = Column(MEDIUMINT(8), nullable=False)
    crt_dateline = Column(INTEGER(10), nullable=False)
    crt_lastpost = Column(INTEGER(11), nullable=False)
    crt_lock = Column(
        TINYINT(4), nullable=False, index=True, server_default=text("'0'")
    )
    crt_img_anidb = Column(VARCHAR(255), nullable=False)
    crt_anidb_id = Column(MEDIUMINT(8), nullable=False)
    crt_ban = Column(TINYINT(3), nullable=False, index=True, server_default=text("'0'"))
    crt_redirect = Column(INTEGER(10), nullable=False, server_default=text("'0'"))
    crt_nsfw = Column(TINYINT(1), nullable=False)


class ChiiCrtCastIndex(Base):
    __tablename__ = "chii_crt_cast_index"

    crt_id = Column(MEDIUMINT(9), primary_key=True, nullable=False)
    prsn_id = Column(MEDIUMINT(9), primary_key=True, nullable=False, index=True)
    subject_id = Column(MEDIUMINT(9), primary_key=True, nullable=False, index=True)
    subject_type_id = Column(
        TINYINT(3), nullable=False, index=True, comment="根据人物归类查询角色，动画，书籍，游戏"
    )
    summary = Column(
        String(255, "utf8_unicode_ci"), nullable=False, comment="幼年，男乱马，女乱马，变身形态，少女形态。。"
    )


class ChiiCrtSubjectIndex(Base):
    __tablename__ = "chii_crt_subject_index"

    crt_id = Column(MEDIUMINT(9), primary_key=True, nullable=False)
    subject_id = Column(MEDIUMINT(9), primary_key=True, nullable=False, index=True)
    subject_type_id = Column(TINYINT(4), nullable=False, index=True)
    crt_type = Column(TINYINT(4), nullable=False, index=True, comment="主角，配角")
    ctr_appear_eps = Column(MEDIUMTEXT, nullable=False, comment="可选，角色出场的的章节")
    crt_order = Column(TINYINT(3), nullable=False)


t_chii_person_alias = Table(
    "chii_person_alias",
    metadata,
    Column("prsn_cat", ENUM("prsn", "crt"), nullable=False),
    Column("prsn_id", MEDIUMINT(9), nullable=False, index=True),
    Column("alias_name", String(255, "utf8_unicode_ci"), nullable=False),
    Column("alias_type", TINYINT(4), nullable=False),
    Column("alias_key", String(10, "utf8_unicode_ci"), nullable=False),
    Index("prsn_cat", "prsn_cat", "prsn_id"),
)


class ChiiPersonCollect(Base):
    __tablename__ = "chii_person_collects"
    __table_args__ = (
        Index("prsn_clt_cat", "prsn_clt_cat", "prsn_clt_mid"),
        {"comment": "人物收藏"},
    )

    prsn_clt_id = Column(MEDIUMINT(8), primary_key=True)
    prsn_clt_cat = Column(Enum("prsn", "crt"), nullable=False)
    prsn_clt_mid = Column(MEDIUMINT(8), nullable=False, index=True)
    prsn_clt_uid = Column(MEDIUMINT(8), nullable=False, index=True)
    prsn_clt_dateline = Column(INTEGER(10), nullable=False)


class ChiiPersonCsIndex(Base):
    __tablename__ = "chii_person_cs_index"
    __table_args__ = {"comment": "subjects' credits/creator & staff (c&s)index"}

    prsn_type = Column(ENUM("prsn", "crt"), primary_key=True, nullable=False)
    prsn_id = Column(MEDIUMINT(9), primary_key=True, nullable=False, index=True)
    prsn_position = Column(
        SMALLINT(5), primary_key=True, nullable=False, index=True, comment="监督，原案，脚本,.."
    )
    subject_id = Column(MEDIUMINT(9), primary_key=True, nullable=False, index=True)
    subject_type_id = Column(TINYINT(4), nullable=False, index=True)
    summary = Column(MEDIUMTEXT, nullable=False)
    prsn_appear_eps = Column(MEDIUMTEXT, nullable=False, comment="可选，人物参与的章节")


class ChiiPersonField(Base):
    __tablename__ = "chii_person_fields"

    prsn_cat = Column(ENUM("prsn", "crt"), primary_key=True, nullable=False)
    prsn_id = Column(INTEGER(8), primary_key=True, nullable=False, index=True)
    gender = Column(TINYINT(4), nullable=False)
    bloodtype = Column(TINYINT(4), nullable=False)
    birth_year = Column(YEAR(4), nullable=False)
    birth_mon = Column(TINYINT(2), nullable=False)
    birth_day = Column(TINYINT(2), nullable=False)


# deprecated table
t_chii_person_relationship = Table(
    "chii_person_relationship",
    metadata,
    Column("prsn_type", ENUM("prsn", "crt"), nullable=False),
    Column("prsn_id", MEDIUMINT(9), nullable=False),
    Column("relat_prsn_type", ENUM("prsn", "crt"), nullable=False),
    Column("relat_prsn_id", MEDIUMINT(9), nullable=False),
    Column("relat_type", SMALLINT(6), nullable=False, comment="任职于，从属,聘用,嫁给，"),
)


class ChiiPerson(Base):
    __tablename__ = "chii_persons"
    __table_args__ = {"comment": "（现实）人物表"}

    prsn_id = Column(MEDIUMINT(8), primary_key=True)
    prsn_name = Column(String(255, "utf8_unicode_ci"), nullable=False)
    prsn_type = Column(TINYINT(4), nullable=False, index=True, comment="个人，公司，组合")
    prsn_infobox = Column(MEDIUMTEXT, nullable=False)
    prsn_producer = Column(TINYINT(1), nullable=False, index=True)
    prsn_mangaka = Column(TINYINT(1), nullable=False, index=True)
    prsn_artist = Column(TINYINT(1), nullable=False, index=True)
    prsn_seiyu = Column(TINYINT(1), nullable=False, index=True)
    prsn_writer = Column(
        TINYINT(4), nullable=False, index=True, server_default=text("'0'"), comment="作家"
    )
    prsn_illustrator = Column(
        TINYINT(4), nullable=False, index=True, server_default=text("'0'"), comment="绘师"
    )
    prsn_actor = Column(TINYINT(1), nullable=False, index=True, comment="演员")
    prsn_summary = Column(MEDIUMTEXT, nullable=False)
    prsn_img = Column(String(255, "utf8_unicode_ci"), nullable=False)
    prsn_img_anidb = Column(VARCHAR(255), nullable=False)
    prsn_comment = Column(MEDIUMINT(9), nullable=False)
    prsn_collects = Column(MEDIUMINT(8), nullable=False)
    prsn_dateline = Column(INTEGER(10), nullable=False)
    prsn_lastpost = Column(INTEGER(11), nullable=False)
    prsn_lock = Column(TINYINT(4), nullable=False, index=True)
    prsn_anidb_id = Column(MEDIUMINT(8), nullable=False)
    prsn_ban = Column(
        TINYINT(3), nullable=False, index=True, server_default=text("'0'")
    )
    prsn_redirect = Column(INTEGER(10), nullable=False, server_default=text("'0'"))
    prsn_nsfw = Column(TINYINT(1), nullable=False)

    def role(self):
        return {
            "producer": self.prsn_producer,
            "mangaka": self.prsn_mangaka,
            "artist": self.prsn_artist,
            "seiyu": self.prsn_seiyu,
            "writer": self.prsn_writer,
            "illustrator": self.prsn_illustrator,
            "actor": self.prsn_actor,
        }