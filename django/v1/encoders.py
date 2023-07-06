from django.core.serializers.json import DjangoJSONEncoder
from django.forms import model_to_dict
from django.db.models import Model   

class ExtendedEncoder(DjangoJSONEncoder):
    def default(self, o):
        if isinstance(o, Model):
            oo = model_to_dict(o)
            for field in o._meta.fields:
                val = getattr(o, field.name)
                if isinstance(val, Model):
                    oo[field.name] = self.default(val)
            return oo
        elif isinstance(o, list):
            return [self.default(item) for item in o]
        return super().default(o)

def ModelToDict(o):
        if isinstance(o, Model):
            oo = model_to_dict(o)
            for field in o._meta.fields:
                val = getattr(o, field.name)
                if isinstance(val, Model):
                    oo[field.name] = ModelToDict(val)
            return oo
        elif isinstance(o, list):
            return [ModelToDict(item) for item in o]
        else:
            return o