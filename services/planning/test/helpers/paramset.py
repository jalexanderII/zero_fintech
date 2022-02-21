import attr


@attr.s(auto_attribs=True, kw_only=True)
class ParamSet:
    """
    Use as base class for sets for parameters for @pytest.mark.parametrize
    when you find yourself passing too many parameters positionally.
    """

    id: str

    @property
    def __name__(self) -> str:  # indicate the id to pytest
        return self.id
