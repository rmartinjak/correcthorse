include config.mk

SRCDIR = src
OBJDIR = obj

_OBJ = correcthorse wlist
OBJ = $(addprefix $(OBJDIR)/,$(addsuffix .o,$(_OBJ)))

all : $(OBJDIR) $(OBJ) correcthorse

$(OBJDIR) :
	@mkdir $(OBJDIR)

$(OBJDIR)/%.o : $(SRCDIR)/%.c
	@echo CC -c $<
	@$(CC) $(CPPFLAGS) $(CFLAGS) -c -o $@ $<

correcthorse : $(OBJ) $(OBJDIR)/correcthorse.o
	@$(CC) $(CPPFLAGS) $(CFLAGS) -o $@ $^ $(LDFLAGS) $(LIBS)

install : all
	@echo installing executable file to ${DESTDIR}${PREFIX}/bin
	@install -d ${DESTDIR}${PREFIX}/bin
	@install -m 755 correcthorse ${DESTDIR}${PREFIX}/bin
	@echo installing wordlists to ${DESTDIR}${PREFIX}/share/correcthorse
	@install -d ${DESTDIR}${PREFIX}/share/correcthorse
	@install -m 644 wordlists/* ${DESTDIR}${PREFIX}/share/correcthorse

clean :
	@echo cleaning
	@rm -rf $(OBJDIR)

.PHONY: clean install
